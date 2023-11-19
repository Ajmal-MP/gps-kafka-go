// server.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Location struct to represent location data
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Server struct to hold request count
type Server struct {
	mu    sync.Mutex
	Count int
}

func (s *Server) handleLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loc Location
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loc)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Increment request count
	s.mu.Lock()
	s.Count++
	count := s.Count
	s.mu.Unlock()

	fmt.Printf("Received request #%d\n", count)

	w.WriteHeader(http.StatusOK)
}

func main() {
	server := &Server{}

	http.HandleFunc("/location", server.handleLocation)

	port := 8081
	fmt.Printf("Server running on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
