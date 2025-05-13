package pool

import (
	"errors"

	"github.com/minhanocaike/smart-balancer/internal/backend"
)

type Pool interface {
	AddBackend(back backend.Backend)
	RemoveBackend(back backend.Backend) error
	ListBackend() []backend.Backend
	Scan(back backend.Backend) error
}

type pool struct {
	backendPool map[string]backend.Backend
}

func NewPool() (Pool, error) {
	return &pool{
		backendPool: make(map[string]backend.Backend),
	}, nil
}

func (p *pool) AddBackend(back backend.Backend) {
	key := getKey(back)
	p.backendPool[key] = back
}

func (p *pool) RemoveBackend(back backend.Backend) error {
	key := getKey(back)
	if _, ok := p.backendPool[key]; !ok {
		return errors.New("host and port not found")
	}

	delete(p.backendPool, key)
	return nil
}

func (p *pool) ListBackend() []backend.Backend {
	var backendList []backend.Backend
	for _, b := range p.backendPool {
		backendList = append(backendList, b)
	}

	return backendList
}

func (p *pool) Scan(back backend.Backend) error {
	return nil
}

func getKey(back backend.Backend) string {
	key := back.Addr.Host + ":" + back.Addr.Port
	return key
}
