package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/server"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {

	storage, err := db.NewStorage(nil)
	handleErr(err)

	newProxy, err := server.NewProxy(context.Background(), storage)
	handleErr(err)

	newProxy.Dispatch(req.Context(), *req)

	fmt.Fprintf(w, "Ola, tudo bem?")
}

func main() {

	srv := http.Server{
		Addr: ":3000",
	}

	http.HandleFunc("/hello", HelloHandler)

	srv.ListenAndServe()
}
