package pool

import (
	"errors"

	"github.com/ckminhano/smart-balancer/internal/backend"
)

type Pool struct {
	BackendPool map[string]backend.Backend
}

func NewPool() (*Pool, error) {
	return &Pool{
		BackendPool: make(map[string]backend.Backend),
	}, nil
}

func (p *Pool) AddBackend(back backend.Backend) {
	key := getKey(back)
	p.BackendPool[key] = back
}

func (p *Pool) RemoveBackend(back backend.Backend) error {
	key := getKey(back)
	if _, ok := p.BackendPool[key]; !ok {
		return errors.New("host and port not found")
	}

	delete(p.BackendPool, key)
	return nil
}

func (p *Pool) ListBackend() []backend.Backend {
	var backendList []backend.Backend
	for _, b := range p.BackendPool {
		backendList = append(backendList, b)
	}

	return backendList
}

func (p *Pool) Scan(back backend.Backend) error {
	return nil
}

func getKey(back backend.Backend) string {
	key := back.Addr.Host + ":" + back.Addr.Port
	return key
}
