package db

import (
	"database/sql"
	"errors"

	"github.com/ckminhano/smart-balancer/internal/router"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) (*Storage, error) {
	if db == nil {
		return &Storage{}, errors.New("db client cannot be nil")
	}

	return &Storage{
		db: db,
	}, nil
}

func (st *Storage) AddRoute(route *router.Route) error {
	// TODO: Implement me

	return nil
}

func (st *Storage) RemoveRoute() error {
	// TODO: Implement me

	return nil
}

func (st *Storage) List() error {
	// TODO: Implement me

	return nil
}
