package pool

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/ckminhano/smart-balancer/internal/backend"
)

type Balance interface {
	Select() *backend.Backend
}

type LeastConnection struct {
	Least  *backend.Backend
	Conns  uint32
	Logger *slog.Logger

	mu sync.Mutex
}

// SelectBackend calls the balancer algorithm to select the backend
func (lc *LeastConnection) Select() (*backend.Backend, error) {
	//TODO: Implement the interface for Balance
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var selected *backend.Backend

	if selected == nil {
		return nil, errors.New("not exists available backends in the pool")
	}

	lc.Logger.Info("selected backend", "id", selected.Id, "host", selected.Addr.Host, "conns", selected.TotalConn())
	return selected, nil
}
