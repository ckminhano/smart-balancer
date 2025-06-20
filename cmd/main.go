package main

import (
	"log"
	"log/slog"

	"github.com/ckminhano/smart-balancer/internal/server"
	"github.com/ckminhano/smart-balancer/internal/spec"
	"github.com/ckminhano/smart-balancer/internal/storage"
)

func main() {
	proxy := load()
	addr := ":3000"

	s, err := server.NewServer(addr, proxy)
	if err != nil {
		log.Panicf("error to create a new server: %v", err)
	}

	err = s.Serve()
	if err != nil {
		log.Panicf("error to start server: %v", err)
	}
}

func load() *server.Proxy {
	// TODO: Load path from flag
	path := "config"
	logger := slog.Default()

	s, err := spec.LoadSpec(path)
	if err != nil {
		log.Panicf("error to load specification: %v", err)
	}

	storage, err := storage.NewStorage(path, s, logger)
	if err != nil {
		log.Panicf("error to load storage: %v", err)
	}

	proxy, err := server.NewProxy(storage)
	if err != nil {
		log.Panicf("error to create a new proxy: %v", err)
	}

	return proxy
}
