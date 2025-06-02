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

func MainHandler(w http.ResponseWriter, req *http.Request) {

	storage, err := db.NewStorage(nil)
	handleErr(err)

	newProxy, err := server.NewProxy(context.Background(), storage)
	handleErr(err)

	res := make(chan http.Response)
	_ = newProxy.Dispatch(req.Context(), res, req)

	_, _ = fmt.Fprintf(w, "Test MainHandler")
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
