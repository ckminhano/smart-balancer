package backend

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
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

	b.Addr.Port = buildPort(b.Addr.Host)

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

func (back *Backend) Invoke(ctx context.Context, res chan<- http.Response, req *http.Request) error {
	headers := req.Header.Clone()
	path := req.URL.Path
	body := req.Body
	protocol := req.Proto

	newURL := url.URL{
		Host:   back.Addr.Host,
		Path:   path,
		Scheme: req.URL.Scheme,
	}

	log.Println(newURL)

	newReq := http.Request{
		Method: req.Method,
		URL:    &newURL,
		Body:   body,
		Header: headers,
		Proto:  protocol,
	}

	client := http.Client{}

	resBackend, err := client.Do(&newReq)
	if err != nil {
		return err
	}

	res <- *resBackend

	return nil
}

// HealthCheck checks if the backend is healthy by sending a GET request to the health path.
// Any response not equal to 200 returns an error
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
