package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dissmaestro/go1/cmd/storage"
)

type Server struct {
	storage *storage.MemStorage
}

func NewServer() *Server {
	return &Server{storage: storage.NewMemStorage()}
}

func (s *Server) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}
	m := storage.Metric{Type: parts[2], Name: parts[3], Value: parts[4]}

	var value interface{}
	var err error

	if m.Type == storage.GaugeMetric {
		value, err = strconv.ParseFloat(m.Value.(string), 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
	} else if m.Type == storage.CounterMetric {
		value, err = strconv.ParseInt(m.Value.(string), 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	m.Value = value
	err = s.storage.UpdateMerics(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}
	data, exist := s.storage.GetMetric(parts[2])
	if !exist {
		http.Error(w, "This metric not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Metric found: %+v\n", data)
}

func main() {
	server := NewServer()
	http.HandleFunc("/update/", server.UpdateMetricHandler)
	http.HandleFunc("/get/", server.GetMetricHandler)
	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}
