package db

import (
	"database/sql"
	"errors"

	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/ckminhano/smart-balancer/internal/route"
	"github.com/ckminhano/smart-balancer/pkg/id"
)

type Storage struct {
	db *sql.Conn
}

func NewStorage(db *sql.DB) (*Storage, error) {
	// if db == nil {
	// 	return &Storage{}, errors.New("db client cannot be nil")
	// }

	return &Storage{
		db: db,
	}, nil
}

func (st *Storage) GetTarget(host string) (*pool.Pool, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}

	// TODO: Implement me
	// Search a target by the specified host source

	return nil, nil
}

func (st *Storage) AddRoute(route *route.Route) error {
	// TODO: Implement me

	return nil
}

func (st *Storage) RemoveRoute(routeId id.Id) error {
	// TODO: Implement me
	return nil
}

func (st *Storage) List() ([]*route.Route, error) {
	// TODO: Implement me
	routes := make([]*route.Route, 0)

	return routes, nil
}
