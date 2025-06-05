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
}

func NewProxy(ctx context.Context, storage *db.Storage) (*Proxy, error) {
	return &Proxy{
		Ctx: ctx,
		Db:  storage,
	}, nil
}

func (p *Proxy) Forward(ctx context.Context, req *http.Request) (*http.Response, error) {
	host := req.Host

	if host == "" {
		return nil, errors.New("could not identify host, check host header value")
	}

	targetPool, err := p.Db.GetTarget(host)
	if err != nil {
		return nil, err
	}

	resCh := make(chan *http.Response, 1)
	errCh := make(chan error, 1)

	go func() {
		err = targetPool.Dispatch(ctx, resCh, req)
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case resp := <-resCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
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
