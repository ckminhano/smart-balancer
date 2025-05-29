package db

import (
	"database/sql"
	"errors"
)

type Router interface {
	AddSource()
	RemoveSource()
	GetTarget()
}

type routerSQL struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) (Router, error) {
	if db == nil {
		return &routerSQL{}, errors.New("db client cannot be nil")
	}

	return &routerSQL{
		db: db,
	}, nil
}

func (r *routerSQL) AddSource() {

}

func (r *routerSQL) RemoveSource() {
	// TODO: Implement me
}

func (r *routerSQL) GetTarget() {
	// TODO: Implement me
}
