package backend

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	Addr *Address
	URL  *string

	// Time in seconds of the timeout duration
	Timeout int16

	// HealthPath is the path to verify the backend health
	HealthPath *string
	client     *http.Client
}

// NewBackend creates a new backend instance with the given address.
func NewBackend(opts ...Option) (*Backend, error) {
	b := &Backend{
		Timeout: 30,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.Addr == nil {
		return nil, errors.New("address cannot be nil, use WithAddress to configure host and port")
	}

	defaultUrl := buildAddress(b.Addr.Host, buildPort(b.Addr.Port), "")

	b.URL = &defaultUrl
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

// HealthCheck checks if the backend is healthy by sending a GET request to the health path.
func (b *Backend) HealthCheck() (int, error) {
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

func buildAddress(host, port, path string) string {
	return fmt.Sprintf("http://%s:%s/%s", host, port, path)
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

func buildPort(port string) string {
	if port == "" {
		return "80"
	}

	return port
}
