package pool

import (
	"context"
	"errors"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/backend"
)

type Host string

type Pool struct {
	BackendPool map[string]*backend.Backend
}

func NewPool() (*Pool, error) {
	return &Pool{
		BackendPool: make(map[string]*backend.Backend),
	}, nil
}

// Receive the request
func (p *Pool) Dispatch(ctx context.Context, res chan<- *http.Response, req *http.Request) error {
	if req == nil {
		return errors.New("http request cannot be nil")
	}

	dest := p.pickBackend(ctx)

	err := dest.Invoke(ctx, res, req)
	if err != nil {
		return err
	}

	return nil
}

// Call the balancer algorithm to select the backend
func (p *Pool) pickBackend(ctx context.Context) *backend.Backend {
	// TODO: Implement me

	return nil
}

func (p *Pool) AddBackend(back *backend.Backend) {
	key := getKey(*back)
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

func (p *Pool) ListBackend() []*backend.Backend {
	var backendList []*backend.Backend
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
