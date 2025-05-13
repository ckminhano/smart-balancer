package pool_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/minhanocaike/smart-balancer/internal/backend"
	"github.com/minhanocaike/smart-balancer/internal/pool"
	"github.com/stretchr/testify/assert"
)

func Main(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestPool_AddBackend(t *testing.T) {
	addrOption := backend.WithAddr(backend.Address{
		Host: "localhost",
		Port: "3000",
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

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			newPool, err := pool.NewPool()
			assert.Equal(t, test.expectedErr, err)

			newPool.AddBackend(test.back)

			poolBackend := newPool.ListBackend()
			fmt.Println(poolBackend)
			assert.Equal(t, test.totalBackend, len(poolBackend))
		})
	}
}

func TestPool_RemoveBackend(t *testing.T) {
	addrOption := backend.WithAddr(backend.Address{
		Host: "localhost",
		Port: "3000",
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

	newPool, err := pool.NewPool()
	assert.Equal(t, testCaseRemovePool.expectedErr, err)

	newPool.AddBackend(testCaseRemovePool.back)

	err = newPool.RemoveBackend(testCaseRemovePool.back)
	assert.Equal(t, testCaseRemovePool.expectedErr, err)

	poolBackend := newPool.ListBackend()
	assert.Equal(t, testCaseRemovePool.totalBackend, len(poolBackend))

	testCaseInvalidPool := struct {
		back        backend.Backend
		expectedErr error
	}{
		back:        *backEmpty,
		expectedErr: errors.New("host and port not found"),
	}

	_, err = pool.NewPool()
	assert.NoError(t, err)

	err = newPool.RemoveBackend(testCaseInvalidPool.back)
	assert.Equal(t, testCaseInvalidPool.expectedErr, err)
}
