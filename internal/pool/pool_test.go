package pool_test

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/ckminhano/smart-balancer/internal/pool"
	"github.com/stretchr/testify/assert"
)

func Main(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestPool_AddBackend(t *testing.T) {
	addrOption := backend.WithAddr(backend.Address{
		Host: "localhost:9000",
	})

	back, err := backend.NewBackend(addrOption)
	assert.NoError(t, err)

	testCase := []struct {
		name         string
		back         backend.Backend
		totalBackend int
		expectedErr  error
	}{
		{
			name:         "Valid Pool",
			back:         *back,
			totalBackend: 1,
			expectedErr:  nil,
		},
	}

	logger := slog.Default()
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			newPool, err := pool.NewPool(logger)
			assert.Equal(t, test.expectedErr, err)

			err = newPool.Add(&test.back)
			assert.Equal(t, test.expectedErr, err)

			poolBackend := newPool.List()
			assert.Equal(t, test.totalBackend, len(poolBackend))
		})
	}
}

func TestPool_RemoveBackend(t *testing.T) {
	t.Skip("Skip remove for now")
	addrOption := backend.WithAddr(backend.Address{
		Host: "localhost:9000",
	})

	back, err := backend.NewBackend(addrOption)
	assert.NoError(t, err)

	addrEmtpyOption := backend.WithAddr(backend.Address{
		Host: "",
		Port: "",
	})

	backEmpty, err := backend.NewBackend(addrEmtpyOption)
	assert.NoError(t, err)

	testCaseRemovePool := struct {
		back         backend.Backend
		totalBackend int
		expectedErr  error
	}{
		back:         *back,
		totalBackend: 0,
		expectedErr:  nil,
	}

	logger := slog.Default()
	newPool, err := pool.NewPool(logger)
	assert.Equal(t, testCaseRemovePool.expectedErr, err)

	err = newPool.Add(&testCaseRemovePool.back)
	assert.Equal(t, testCaseRemovePool.expectedErr, err)

	err = newPool.Remove(testCaseRemovePool.back)
	assert.Equal(t, testCaseRemovePool.expectedErr, err)

	poolBackend := newPool.List()
	assert.Equal(t, testCaseRemovePool.totalBackend, len(poolBackend))

	testCaseInvalidPool := struct {
		back        backend.Backend
		expectedErr error
	}{
		back:        *backEmpty,
		expectedErr: errors.New("backend id not found in the pool"),
	}

	_, err = pool.NewPool(logger)
	assert.NoError(t, err)

	err = newPool.Remove(testCaseInvalidPool.back)
	assert.Equal(t, testCaseInvalidPool.expectedErr, err)
}

func TestPool_SelectBackend(t *testing.T) {
	// TODO: Improve testCase for range loop

	addrOption1 := backend.WithAddr(backend.Address{
		Host: "localhost:9000",
	})
	back1, err := backend.NewBackend(addrOption1)
	assert.NoError(t, err)
	back1.Connections = 10

	addrOption2 := backend.WithAddr(backend.Address{
		Host: "localhost:9001",
	})
	back2, err := backend.NewBackend(addrOption2)
	assert.NoError(t, err)
	back2.Connections = 5

	addrOption3 := backend.WithAddr(backend.Address{
		Host: "localhost:9002",
	})
	back3, err := backend.NewBackend(addrOption3)
	assert.NoError(t, err)
	back3.Connections = 15

	logger := slog.Default()
	p, err := pool.NewPool(logger)
	assert.NoError(t, err)

	testCaseSelect := struct {
		testPool        *pool.Pool
		selectedBackend *backend.Backend
	}{
		testPool:        p,
		selectedBackend: back2,
	}

	backs := []*backend.Backend{back1, back2, back3}

	for _, back := range backs {
		err = testCaseSelect.testPool.Add(back)
		assert.NoError(t, err)
	}

	selected, err := testCaseSelect.testPool.Best()
	assert.NoError(t, err)

	fmt.Println(selected.Addr.Host)
	for i, b := range testCaseSelect.testPool.List() {
		fmt.Printf("Index: %v, Host: %v\n", i, b.Addr.Host)
	}

	assert.Equal(t, testCaseSelect.selectedBackend, selected)
}
