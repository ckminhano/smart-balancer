package server

import (
	"context"
	"errors"

	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/router"
	"github.com/ckminhano/smart-balancer/pkg/id"
	"github.com/google/uuid"
)

type Proxy struct {
	Ctx context.Context
	Db  db.Storage
}

func NewProxy(ctx context.Context, storage db.Storage) (*Proxy, error) {
	return &Proxy{
		Ctx: ctx,
		Db:  storage,
	}, nil
}

func (p *Proxy) AddRoute(route *router.Route) error {
	if route.Id.UUID() == uuid.Nil {
		route.Id = id.NewId()
	}

	if route.Source == "" {
		return errors.New("route source cannot be empty")
	}

	err := p.Db.AddRoute(route)
	if err != nil {
		return err
	}

	return nil
}
