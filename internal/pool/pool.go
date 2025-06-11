package pool

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"sync"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/google/uuid"
)

type Host string

type Pool struct {
	Backends []*backend.Backend
	Logger   *slog.Logger
	mu       sync.RWMutex

	leastIdxConnection int32
}

func NewPool(logger *slog.Logger) (*Pool, error) {
	return &Pool{
		Backends: make([]*backend.Backend, 0),
		Logger:   logger,
	}, nil
}

// Dispatch receives and start the request to the backend server.
func (p *Pool) Dispatch(ctx context.Context, res chan<- *http.Response, req *http.Request) error {
	if req == nil {
		return errors.New("http request cannot be nil")
	}

	dest, err := p.Select()
	if err != nil {
		return err
	}

	err = dest.Invoke(ctx, res, req)
	if err != nil {
		return err
	}

	return nil
}

// SelectBackend calls the balancer algorithm to select the backend
func (p *Pool) Select() (*backend.Backend, error) {
	maxConns := int32(math.MaxInt32)

	var selected *backend.Backend
	for _, b := range p.Backends {
		if b.TotalConn() <= maxConns {
			selected = b
			maxConns = b.TotalConn()
		}
	}

	if selected == nil {
		p.Logger.Error("backends not available", "pool", p.Backends)
		return nil, errors.New("not exists available backends in the pool")
	}

	p.Logger.Info("selected backend", "id", selected.Id, "host", selected.Addr.Host, "conns", selected.TotalConn())
	return selected, nil
}

func (p *Pool) Add(back *backend.Backend) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, b := range p.Backends {
		if b.Addr.Host == back.Addr.Host {
			return errors.New("backend with this address already exists in this pool")
		}
	}

	p.Backends = append(p.Backends, back)

	return nil
}

func (p *Pool) Remove(back backend.Backend) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if back.Id == nil || back.Id.UUID() == uuid.Nil {
		return errors.New("invalid backend id")
	}

	idx := -1
	for i, b := range p.Backends {
		if b.Id == back.Id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("backend id not found in the pool")
	}

	p.Backends = slices.Delete(p.Backends, idx, idx+1)

	return nil
}

// List safes list the backends
func (p *Pool) List() []*backend.Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.Backends
}

// Scan checks periodically the health of the backends in the pool.
// In case of a backend failure, it should remove the backend from the pool.
func (p *Pool) Scan(back backend.Backend) error {
	return nil
}
