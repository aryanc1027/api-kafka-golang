package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL    = "http://localhost:8080"
	numStreams = 1000
	apiKey     = "your-api-key-here"
)

func main() {
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < numStreams; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			streamID := startStream()
			sendData(streamID)
		}()
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("Benchmark completed in %v\n", duration)
	fmt.Printf("Average time per stream: %v\n", duration/numStreams)
}

func startStream() string {
	resp, err := http.Post(baseURL+"/api/stream/start", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["stream_id"]
}

func sendData(streamID string) {
	data := map[string]string{"message": "test data"}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/stream/%s/send", baseURL, streamID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}