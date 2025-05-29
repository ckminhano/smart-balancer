package router

import (
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
