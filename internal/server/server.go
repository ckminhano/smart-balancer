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
	Proxy *Proxy
	srv   http.Server
	ctx   context.Context
}

func NewServer(addr string, proxy *Proxy) (*Server, error) {
	if proxy.Db == nil {
		return nil, errors.New("proxy db cannot be empty")
	}

	return &Server{
		Proxy: proxy,
		srv: http.Server{
			Addr: addr,
		},
	}, nil
}

func (s *Server) ProxyHandler(w http.ResponseWriter, req *http.Request) {
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

	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println("error copying body to client: ", err.Error())
	}
}

func (s *Server) Serve() error {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	s.srv.Handler = http.HandlerFunc(s.ProxyHandler)
	s.srv.ReadTimeout = 10 * time.Second
	s.srv.WriteTimeout = 10 * time.Second

	go func() {
		log.Println("starting proxy server on ", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-sigs
	log.Println("shutting down proxy server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server gracefully stopped.")

	return nil
}
