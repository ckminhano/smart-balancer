package route

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/pkg/id"
)

type Route struct {
	Id   *id.Id
	Name *string

	// Source uses a string to represent the host
	Source string
	Target *pool.Pool
}

func NewRoute(name string, target *pool.Pool, src string) (*Route, error) {
	if src == "" {
		return nil, errors.New("route source cannot be empty")
	}
	if target == nil {
		return nil, errors.New("target cannot be nil")
	}

	return &Route{
		Id:     id.NewId(),
		Name:   &name,
		Source: src,
		Target: target,
	}, nil
}
