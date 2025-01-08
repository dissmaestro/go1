package main

import (
	"fmt"
	"net/http"

	Storage "github.com/dissmaestro/go1/storage"
)

type Server struct {
	storage *Storage.MemStorage
}

func run(mux *http.ServeMux) error {
	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		return fmt.Errorf("Failed to listen and serve: %w", err)
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	if err := run(mux); err != nil {
		panic(err)
	}
}
