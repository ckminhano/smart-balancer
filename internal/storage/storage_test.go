package storage_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/ckminhano/smart-balancer/internal/spec"
	"github.com/ckminhano/smart-balancer/internal/storage"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test_StorageLoad(t *testing.T) {
	path := "config"

	s, err := spec.LoadSpec(path)
	if err != nil {
		t.Errorf("error to load specification: %v", err)
	}

	logger := slog.Default()
	st, err := storage.NewStorage(path, s, logger)
	if err != nil {
		t.Errorf("error to create a new storage: %v", err)
	}

	routes := st.List()

	// TODO: Create a assert function to check if the routes are loaded correctly
	for _, r := range routes {
		fmt.Printf("route name: %s\n", *r.Name)
		fmt.Printf("route origin: %s\n", r.Origin)
		fmt.Printf("route id: %s\n", r.Id)

		for _, b := range r.Target.List() {
			fmt.Printf("backend host: %s\n", b.Addr.Host)
			fmt.Printf("backend health path: %s\n", *b.HealthPath)
		}
	}
}
