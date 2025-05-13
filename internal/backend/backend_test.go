package backend_test

import (
	"os"
	"testing"

	"github.com/minhanocaike/smart-balancer/internal/backend"
	"github.com/stretchr/testify/assert"
)

func Main(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestNewBackend(t *testing.T) {
	testCase := []struct {
		name         string
		host         string
		port         string
		healthPath   string
		responseCode int
		expectedErr  error
	}{
		{
			name:         "Valid Backend",
			host:         "localhost",
			port:         "3000",
			healthPath:   "/",
			responseCode: 200,
			expectedErr:  nil,
		},
		{
			name:         "Invalid Path",
			host:         "localhost",
			port:         "3000",
			healthPath:   "",
			responseCode: 0,
			expectedErr:  backend.ErrorPath,
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			addrOption := backend.WithAddr(backend.Address{
				Host: test.host,
				Port: test.port,
			})
			healthPathOption := backend.WithHealthPath(test.healthPath)

			backend, err := backend.NewBackend(addrOption, healthPathOption)
			assert.NoError(t, err)

			status, err := backend.HealthCheck()
			assert.Equal(t, test.expectedErr, err)

			assert.Equal(t, test.responseCode, status)
		})
	}
}
