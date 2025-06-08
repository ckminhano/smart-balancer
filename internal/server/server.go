package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http.Server

	Proxy *Proxy
	ctx   context.Context
}

func NewServer(addr string, proxy *Proxy) (*Server, error) {
	if proxy.Storage == nil {
		return nil, errors.New("proxy storage cannot be empty")
	}

	return &Server{
		Proxy: proxy,
		Server: http.Server{
			Addr: addr,
		},
	}, nil
}

func (s *Server) ProxyHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("received request: ", req.Method, req.URL.Path)
	res, err := s.Proxy.Forward(s.ctx, req)
	if err != nil {
		log.Printf("proxy error: %v\n", err)
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Printf("error to close body: %v", err)
		}
	}()

	// Iterate over the header hListValue because the header can have multiple values with a slice of strings
	// and hValue is a single value of the header
	for hName, hListValue := range res.Header {
		for _, hValue := range hListValue {
			w.Header().Add(hName, hValue)
		}
	}

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println("error copying body to client: ", err.Error())
	}
	log.Println("executed request: ", req.Method, req.URL.Path, "with status code: ", res.StatusCode)
}

func (s *Server) Serve() error {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	s.Handler = http.HandlerFunc(s.ProxyHandler)
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second

	go func() {
		log.Println("starting proxy server on ", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-sigs
	log.Println("shutting down proxy server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server gracefully stopped.")

	return nil
}
