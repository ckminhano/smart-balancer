package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/route"
)

type Proxy struct {
	Ctx context.Context
	Db  *db.Storage
}

func NewProxy(ctx context.Context, storage *db.Storage) (*Proxy, error) {
	return &Proxy{
		Ctx: ctx,
		Db:  storage,
	}, nil
}

func (p *Proxy) Forward(ctx context.Context, res chan<- http.Response, req *http.Request) error {
	host := req.Host

	if host == "" {
		return errors.New("could not identify host, check host header value")
	}

	fmt.Println("path: ", req.URL.Path)

	targetPool, err := p.Db.GetTarget(host)
	if err != nil {
		return err
	}

	err = targetPool.Dispatch(ctx, res, req)
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
