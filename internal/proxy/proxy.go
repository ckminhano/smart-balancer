package proxy

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/pool"
)

type Route struct {
	id     int
	name   *string
	Source string
	Target *pool.Pool
}

type Proxy struct {
	Router map[string]*pool.Pool
}

func NewProxy() (*Proxy, error) {
	return &Proxy{}, nil
}

func (p *Proxy) AddRoute(route *Route) error {
	if route.Source == "" || route.Target == nil {
		errors.New("source and target cannot be empty")
	}

}

// generateName returns a new randle name based on the existings route name, but unique
func generateName(route *Route) string {
	backendHost := route.Target
	return route.Source + 
}
