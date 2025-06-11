package backend

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ckminhano/smart-balancer/pkg/id"
	"github.com/google/uuid"
)

type Option func(*Backend)

type Protocol string

const (
	// TODO: Add support for HTTPS
	HTTP  Protocol = "http"
	HTTPS Protocol = "https"
)

// Address represents the address of a backend server.
type Address struct {
	Protocol Protocol
	Host     string
	Port     string
}

type Backend struct {
	*http.Client
	Logger *slog.Logger

	Addr *Address
	URL  *string
	Id   *id.Id

	Timeout    int16
	HealthPath *string

	// Conns is not safe to use, instead get with TotalConn()
	Conns  int32
	Active bool
}

// NewBackend creates a new backend instance with the given address.
func NewBackend(opts ...Option) (*Backend, error) {
	b := &Backend{
		Timeout: 30,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
		Id: id.NewId(),
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.Addr == nil {
		return nil, errors.New("address cannot be nil, use WithAddress to configure host and port")
	}

	if b.Addr.Port == "" {
		b.Addr.Port = buildPort(b.Addr.Host)
	}

	if b.Logger == nil {
		b.Logger = slog.Default()
	}

	return b, nil
}

func WithAddr(addr Address) Option {
	return func(b *Backend) {
		b.Addr = &addr
	}
}

func WithHealthPath(addr string) Option {
	return func(b *Backend) {
		b.HealthPath = &addr
	}
}

func WithTimeout(timeout int16) Option {
	return func(b *Backend) {
		b.Timeout = timeout
	}
}

func WithLogger(logger slog.Logger) Option {
	return func(b *Backend) {
		b.Logger = &logger
	}
}

func (back *Backend) Invoke(ctx context.Context, res chan<- *http.Response, req *http.Request) error {
	// Change request host to backend host
	req.URL.Host = back.Addr.Host

	atomic.AddInt32(&back.Conns, 1)
	defer atomic.AddInt32(&back.Conns, -1)

	back.Logger.Info("backend request", "host", back.Addr.Host, "connections_number", atomic.LoadInt32(&back.Conns))

	backendResp, err := back.Do(req)
	if err != nil {
		return err
	}

	select {
	case res <- backendResp:
		return nil
	case <-ctx.Done():
		err := backendResp.Body.Close()
		if err != nil {
			log.Printf("error to close backend response body: %v", err)
			back.Logger.Info("error to close backend response body", "err", err, "host", back.Addr.Host)
			return err
		}

		return ctx.Err()
	}
}

// HealthCheck checks if the backend is healthy by sending a GET request to the health path.
// Any response not equal to 200 returns an error
func (b *Backend) HealthCheck() (int, error) {
	if b.Id == nil || b.Id.UUID() == uuid.Nil {
		return 0, errors.New("backend id cannot be nil")
	}

	if b.HealthPath == nil {
		return 0, errors.New("health path cannot be empty, use WithHealthPath to configure")
	}

	if exists := strings.HasPrefix(*b.HealthPath, "/"); !exists {
		return 0, errors.New("path must starts with /")
	}

	healthURI := buildURL(*b.URL, *b.HealthPath)

	resp, err := http.Get(healthURI)
	if err != nil || resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("backend not health")
	}

	return resp.StatusCode, nil
}

func (b *Backend) TotalConn() int32 {
	total := atomic.LoadInt32(&b.Conns)

	return total
}

func buildURL(url, path string) string {
	if path == "/" {
		return url
	}

	if cleanURL, exists := strings.CutSuffix(url, "/"); exists {
		return cleanURL + path
	}

	return url + path
}

func buildPort(host string) string {
	if host == "" {
		return "80"
	}

	port := strings.Split(host, ":")[1]

	return port
}
