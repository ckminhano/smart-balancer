package route

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/pkg/id"
)

type Route struct {
	Id   *id.Id
	Name *string

	// Implementar melhoria utilizando um type host
	Source string
	Target *pool.Pool
}

func NewRoute(name string, source string, target *pool.Pool) (*Route, error) {
	if source == "" {
		return nil, errors.New("route source cannot be empty")
	}

	if target == nil {
		return nil, errors.New("target cannot be nil")
	}

	return &Route{
		Id:     id.NewId(),
		Name:   &name,
		Source: source,
		Target: target,
	}, nil
}
