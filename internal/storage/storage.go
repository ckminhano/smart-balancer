package storage

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/internal/route"
	"github.com/ckminhano/smart-balancer/internal/spec"
)

type Storage struct {
	routes        map[string]*route.Route
	specification spec.Specification
}

// NewStorage returns a storage not nil if not err, if path is empty returns an error
// - path: path for the specification file with routes, for spec details see documentation
func NewStorage(path string) (*Storage, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	s, err := spec.LoadSpec(path)
	if err != nil {
		return nil, err
	}

	routes, err := loadRoutes(&s)
	if err != nil {
		return nil, err
	}

	return &Storage{
		routes:        routes,
		specification: s,
	}, nil
}

func (st *Storage) GetTarget(origin string) (*pool.Pool, error) {
	if origin == "" {
		return nil, errors.New("origin cannot be empty")
	}

	if _, ok := st.routes[origin]; !ok {
		return nil, errors.New("origin not found in the routes")
	}
	r := st.routes[origin]

	return r.Target, nil
}

func (st *Storage) List() []*route.Route {
	routes := make([]*route.Route, 0)

	for _, v := range st.routes {
		routes = append(routes, v)
	}

	return routes
}

func loadRoutes(s *spec.Specification) (map[string]*route.Route, error) {
	mapRoute := make(map[string]*route.Route)

	for _, r := range s.Routes {
		loadedPool, err := pool.NewPool()
		if err != nil {
			return nil, err
		}

		for _, loadedBackend := range r.Backends {
			// FIXME: fix protocol to receive from loaded file
			addr := backend.Address{
				Protocol: backend.HTTP,
				Host:     loadedBackend.Host,
			}

			back, err := backend.NewBackend(
				backend.WithAddr(addr),
				backend.WithHealthPath(loadedBackend.Health),
			)
			if err != nil {
				return nil, err
			}

			err = loadedPool.AddBackend(back)
			if err != nil {
				return nil, err
			}
		}

		loadedRoute, err := route.NewRoute(r.Name, loadedPool, r.Origin)
		if err != nil {
			return nil, err
		}

		mapRoute[r.Origin] = loadedRoute
	}

	return mapRoute, nil
}
