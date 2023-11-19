package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Location struct to represent location data
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func generateRandomLocation() Location {
	// For simplicity, generating random coordinates within a specific range.
	// You might need to adjust this based on your requirements.
	return Location{
		Latitude:  rand.Float64()*180 - 90,
		Longitude: rand.Float64()*360 - 180,
	}
}

func sendLocation(loc Location, wg *sync.WaitGroup) {
	defer wg.Done()

	url := "http://localhost:8081/location" // Update with your server URL

	jsonData, err := json.Marshal(loc)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

func main() {
	// Seed the random number generator to get different values each time
	rand.Seed(time.Now().UnixNano())

	// Number of concurrent requests
	numRequests := 10000

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Send location data in parallel
	for i := 0; i < numRequests; i++ {
		loc := generateRandomLocation()
		go sendLocation(loc, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
