package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/route"
)

type Proxy struct {
	Ctx context.Context
	Db  *db.Storage

	// TODO: Verificar se o response vem do server
	Res chan<- *http.Response
}

func NewProxy(ctx context.Context, storage *db.Storage) (*Proxy, error) {
	return &Proxy{
		Ctx: ctx,
		Db:  storage,
		Res: make(chan<- *http.Response),
	}, nil
}

func (p *Proxy) Forward(ctx context.Context, res http.Response, req *http.Request) error {
	host := req.Host
	defer close(p.Res)

	if host == "" {
		return errors.New("could not identify host, check host header value")
	}

	targetPool, err := p.Db.GetTarget(host)
	if err != nil {
		return err
	}

	err = targetPool.Dispatch(ctx, p.Res, req)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) AddRoute(route *route.Route) error {
	err := p.Db.AddRoute(route)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) RemoveRoute(route *route.Route) error {
	err := p.Db.RemoveRoute(*route.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) ListRoutes() ([]*route.Route, error) {
	routes, err := p.Db.List()
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 {
		return nil, errors.New("RoutesNotFound")
	}

	return routes, nil
}
