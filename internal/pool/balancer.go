package pool

import (
	"container/heap"
	"errors"
	"log/slog"

	"github.com/ckminhano/smart-balancer/internal/backend"
)

type Balancer interface {
	// SelectBackend calls the balancer algorithm to select appropriate the backend
	Select() (*backend.Backend, error)

	Insert(back *backend.Backend) error
	List() []*backend.Backend
	Update(idx int) error
	Remove(idx int) error
}

// LeastConnection implements a load balancer that selects the backend with the least number of connections.
// It implements the heap.Interface to allow it to be used with the heap package to manage the priority of backends based on their connection counts.
// It is not safe for concurrent use, so it should be used with a mutex or other synchronization mechanism.
type LeastConnection struct {
	logger *slog.Logger
	B      []*backend.Backend
}

func NewBalancer(logger *slog.Logger) (Balancer, error) {
	if logger == nil {
		return nil, errors.New("backends cannot be empty or nil")
	}

	backs := make([]*backend.Backend, 0)
	lc := &LeastConnection{
		logger: logger,
		B:      backs,
	}

	heap.Init(lc)

	return lc, nil
}

func (lc *LeastConnection) Select() (*backend.Backend, error) {
	if lc.B == nil {
		return nil, errors.New("there a not available backends to select")
	}

	// Takes the first element of the min-heap, that holds a pointer to the backend with least connections
	backLeastConnection := lc.B[0]

	return backLeastConnection, nil
}

func (lc *LeastConnection) Update(idx int) error {
	if len(lc.B) <= idx {
		return errors.New("index not allowed, and out of backends size")
	}

	heap.Fix(lc, idx)
	return nil
}

func (lc *LeastConnection) Insert(back *backend.Backend) error {
	if lc.B == nil || back == nil {
		return errors.New("backends cannot be empty or nil")
	}

	heap.Push(lc, back)
	return nil
}

func (lc *LeastConnection) Remove(idx int) error {
	if len(lc.B) <= idx {
		return errors.New("index not allowed, and out of backends size")
	}

	heap.Remove(lc, idx)
	return nil
}

func (lc *LeastConnection) Push(back any) {
	if lc.B == nil || back == nil {
		lc.logger.Error("pool and backend cannot be emtpy or nil")
		return
	}

	if b, ok := back.(*backend.Backend); ok {
		lc.logger.Info("added new backend", "id", b.Id, "host", b.Addr.Host)
		lc.B = append(lc.B, b)
	}
}

func (lc *LeastConnection) List() []*backend.Backend {
	return lc.B
}

func (lc *LeastConnection) Pop() any {
	if lc.B == nil {
		lc.logger.Error("pool and backend cannot be emtpy or nil")
		return nil
	}

	oldSize := len(lc.B) - 1
	removed := (lc.B)[oldSize]

	lc.B = lc.B[:oldSize-1]

	return removed
}

func (lc *LeastConnection) Len() int {
	if lc.B == nil {
		return 0
	}

	return len(lc.B)
}

func (lc *LeastConnection) Less(i int, j int) bool {
	if lc.B == nil || len(lc.B) <= i || len(lc.B) <= j {
		return false
	}

	backs := lc.B
	if backs[i] == nil || backs[j] == nil {
		return false
	}

	return backs[i].Connections < backs[j].Connections
}

func (lc *LeastConnection) Swap(i int, j int) {
	if lc.B == nil || len(lc.B) <= i || len(lc.B) <= j {
		return
	}

	lc.B[i], lc.B[j] = lc.B[j], lc.B[i]
}
