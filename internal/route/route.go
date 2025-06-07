package route

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/pkg/id"
)

type Route struct {
	Id   *id.Id
	Name *string

	// origin uses a string to represent the fqdn
	Origin string
	Target *pool.Pool
}

// NewRoute create and returns a new NewRoute
// - name: name to identify the route
// - target: pool to route request
// - origin: fqdn or host for the received request and route rule
func NewRoute(name string, target *pool.Pool, origin string) (*Route, error) {
	if origin == "" {
		return nil, errors.New("route source cannot be empty")
	}
	if target == nil {
		return nil, errors.New("target cannot be nil")
	}

	return &Route{
		Id:     id.NewId(),
		Name:   &name,
		Origin: origin,
		Target: target,
	}, nil
}
