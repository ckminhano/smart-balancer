package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ckminhano/smart-balancer/internal/backend"
	"github.com/ckminhano/smart-balancer/internal/db"
	"github.com/ckminhano/smart-balancer/internal/server"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func MainHandler(w http.ResponseWriter, req *http.Request) {
	res, err := CallProxy(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "error to forward backend")
		return
	}
	defer res.Body.Close()

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

	res := make(chan *http.Response)
	defer close(res)
	err = back.Invoke(context.Background(), res, req)
	if err != nil {
		log.Fatal("error to invoke backend: ", err.Error())
	}
}

func CallProxy(req *http.Request) (*http.Response, error) {
	storage, err := db.NewStorage(nil)
	handleErr(err)

	ctx := req.Context()

	testProxy, err := server.NewProxy(context.Background(), storage)
	handleErr(err)

	responseProxy, err := testProxy.Forward(ctx, req)
	if err != nil {
		return nil, err
	}

	return responseProxy, nil
}
