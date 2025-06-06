package main

import (
	"context"
	"log"

	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/server"
)

func main() {
	proxy := load()
	addr := ":3000"

	proxyServer, err := server.NewServer(addr, proxy)
	if err != nil {
		log.Panicf("error to create a new server: %v", err)
	}

	err = proxyServer.Serve()
	if err != nil {
		log.Panicf("error to start server: %v", err)
	}
}

func load() *server.Proxy {
	storage, err := db.NewStorage(nil)
	if err != nil {
		log.Panicf("error to load storage: %v", err)
	}

	proxy, err := server.NewProxy(context.Background(), storage)
	if err != nil {
		log.Panicf("error to create a new proxy: %v", err)
	}

	return proxy
}
