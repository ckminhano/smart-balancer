package spec_test

import (
	"os"
	"testing"

	"github.com/ckminhano/smart-balancer/internal/spec"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestLoad(t *testing.T) {
	path := "config_test"

	routes, err := spec.LoadSpec(path)
	if err != nil {
		t.Errorf("error to load file: %v", err)
	}

	if len(routes.Routes) == 0 {
		t.Error("no routes loaded")
	}
}
