package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/backend"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func MainHandler(w http.ResponseWriter, req *http.Request) {

	// storage, err := db.NewStorage(nil)
	// handleErr(err)

	// newProxy, err := server.NewProxy(context.Background(), storage)
	// handleErr(err)

	// res := make(chan http.Response)
	// _ = newProxy.Forward(req.Context(), res, req)

	CallBackend(req)

	_, _ = fmt.Fprintf(w, "proxy executed")
}

func main() {

	srv := http.Server{
		Addr: ":3000",
	}

	http.HandleFunc("/", MainHandler)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func CallBackend(req *http.Request) {
	addr := backend.WithAddr(backend.Address{
		Protocol: backend.HTTP,
		Host:     "localhost:9000",
	})

	back, err := backend.NewBackend(addr)
	if err != nil {
		log.Fatal("error to create a new backend: ", err.Error())
	}

	res := make(chan http.Response)
	defer close(res)
	err = back.Invoke(context.Background(), res, req)
	if err != nil {
		log.Fatal("error to invoke backend: ", err.Error())
	}
}
