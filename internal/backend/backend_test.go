package backend_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/stretchr/testify/assert"
)

func Main(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestNewBackend(t *testing.T) {
	t.Skip("Skip for now")

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
			host:         "localhost:3000",
			healthPath:   "/",
			responseCode: 200,
			expectedErr:  nil,
		},
		{
			name:         "Invalid Path",
			host:         "localhost:3000",
			healthPath:   "",
			responseCode: 0,
			expectedErr:  backend.ErrorPath,
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			addrOption := backend.WithAddr(backend.Address{
				Host: test.host,
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

func TestBackendInvoke(t *testing.T) {
	headers := http.Header{}
	headers.Add("Host", "localhost:3000")

	testURL := url.URL{
		Scheme: "http",
		Path:   "/smart/balancer",
	}

	reqTest := http.Request{
		Method: "GET",
		Header: headers,
		URL:    &testURL,
		Proto:  "HTTP/1.1",
	}

	addr := backend.WithAddr(backend.Address{
		Protocol: backend.HTTP,
		Host:     "localhost:9000",
	})

	back, err := backend.NewBackend(addr)
	if err != nil {
		assert.Fail(t, "error to create a new backend", err.Error())
	}

	res := make(chan http.Response)
	err = back.Invoke(context.Background(), res, &reqTest)
	if err != nil {
		assert.Fail(t, "error to invoke backend", err.Error())
	}
}
