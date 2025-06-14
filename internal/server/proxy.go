package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/ckminhano/smart-balancer/internal/route"
	"github.com/ckminhano/smart-balancer/internal/storage"
)

type requestParameters struct {
	headers  http.Header
	body     io.ReadCloser
	protocol string
	method   string
	url      url.URL
}

type Proxy struct {
	Ctx     context.Context
	Storage *storage.Storage
}

func NewProxy(storage *storage.Storage) (*Proxy, error) {
	return &Proxy{
		Storage: storage,
	}, nil
}

func (p *Proxy) Forward(ctx context.Context, req *http.Request) (*http.Response, error) {
	host := req.Host

	if host == "" {
		return nil, errors.New("could not identify host, check host header value")
	}

	targetPool, err := p.Storage.GetTarget(host)
	if err != nil {
		return nil, err
	}

	forwardRequest := middleware(req)
	res, err := targetPool.Dispatch(ctx, forwardRequest)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Proxy) ListRoutes() ([]*route.Route, error) {
	routes := p.Storage.List()

	if len(routes) == 0 {
		return nil, errors.New("RoutesNotFound")
	}

	return routes, nil
}

// middleware apply the transformation rules
func middleware(req *http.Request) *http.Request {
	reqParam := requestParameters{
		headers:  req.Header.Clone(),
		body:     req.Body,
		protocol: req.Proto,
		method:   req.Method,
		url: url.URL{
			Path: req.URL.Path,
			// FIXME: This is a temporary solution, we should use the protocol from the request
			Scheme: "http",
		},
	}

	// TODO: Check parser
	// Maybe I can put somen transformations here
	backendRequest := http.Request{
		Method: reqParam.method,
		URL:    &reqParam.url,
		Body:   reqParam.body,
		Header: reqParam.headers,
		Proto:  reqParam.protocol,
	}

	return &backendRequest
}
