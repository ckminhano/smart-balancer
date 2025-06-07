package pool

import (
	"context"
	"errors"
	"net/http"
	"strings"

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

	dest, err := p.pickBackend()
	if err != nil {
		return err
	}

	err = dest.Invoke(ctx, res, req)
	if err != nil {
		return err
	}

	return nil
}

// Call the balancer algorithm to select the backend
func (p *Pool) pickBackend() (*backend.Backend, error) {
	// TODO: Implement me

	// FIXME: Replace the mock back[0] for the right back
	backs := p.ListBackend()

	return backs[0], nil
}

func (p *Pool) AddBackend(back *backend.Backend) error {
	key := getKey(*back)

	if _, ok := p.BackendPool[key]; ok {
		return errors.New("backend already exists in this pool")
	}

	p.BackendPool[key] = back
	return nil
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

// Scan checks periodically the health of the backends in the pool.
// In case of a backend failure, it should remove the backend from the pool.
func (p *Pool) Scan(back backend.Backend) error {
	return nil
}

// getKey generates a unique key for the backend based on its address.
func getKey(back backend.Backend) string {
	key := strings.ToLower(back.Addr.Host + ":" + back.Addr.Port)
	return key
}
