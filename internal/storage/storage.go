package storage

import (
	"errors"
	"log/slog"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/internal/route"
	"github.com/ckminhano/smart-balancer/internal/spec"
)

type Storage struct {
	Logger *slog.Logger

	routes        map[string]*route.Route
	specification spec.Specification
}

// NewStorage returns a storage not nil if not err, if path is empty returns an error
// - path: path for the specification file with routes, for spec details see documentation
func NewStorage(path string, settings spec.Specification, logger *slog.Logger) (*Storage, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	s := &Storage{
		Logger:        logger,
		specification: settings,
	}

	err := s.load()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetTarget search for the origin in the routes, if not found returns an error
// - origin: the fqdn in the host header of the request
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

// load receives the specification file to create the routes
func (st *Storage) load() error {
	mapRoute := make(map[string]*route.Route)

	for _, r := range st.specification.Routes {
		loadedPool, err := pool.NewPool(st.Logger)
		if err != nil {
			return err
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
				return err
			}

			err = loadedPool.Add(back)
			if err != nil {
				return err
			}

			st.Logger.Info("loaded backend", "host", back.Addr.Host, "id", back.Id)
		}

		loadedRoute, err := route.NewRoute(r.Name, loadedPool, r.Origin)
		if err != nil {
			return err
		}
		st.Logger.Info("loaded route", "name", r.Name, "origin", r.Origin)

		mapRoute[r.Origin] = loadedRoute
	}
	st.routes = mapRoute

	return nil
}
