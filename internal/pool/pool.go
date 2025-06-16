package pool

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/google/uuid"
)

type Host string

type Pool struct {
	Logger   *slog.Logger
	balancer Balancer

	mu sync.RWMutex
}

func NewPool(logger *slog.Logger) (*Pool, error) {

	b, err := NewBalancer(logger)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Logger:   logger,
		balancer: b,
	}, nil
}

func (p *Pool) Best() (*backend.Backend, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.balancer.Select()
}

// Dispatch receives and start the request to the backend server.
func (p *Pool) Dispatch(ctx context.Context, req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("http request cannot be nil")
	}

	dest, err := p.balancer.Select()
	if err != nil {
		return nil, err
	}

	res, err := dest.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Pool) Add(back *backend.Backend) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.balancer.Insert(back)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) Remove(back backend.Backend) error {
	// TODO: Implement me
	if back.Id == nil || back.Id.UUID() == uuid.Nil {
		return errors.New("invalid backend id")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.balancer.Remove(0)
	if err != nil {
		return err
	}

	return nil
}

// List safes list the backends
func (p *Pool) List() []*backend.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.balancer.List()
}

// Scan checks periodically the health of the backends in the pool.
// In case of a backend failure, it should remove the backend from the pool.
func (p *Pool) Scan(back backend.Backend) error {
	return nil
}
