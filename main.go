package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

var (
	imageMemory []byte                           // In-memory storage for the screenshot
	mu          sync.Mutex                       // Mutex to manage concurrent access
	clients     = make(map[chan []byte]struct{}) // Map to store SSE clients
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/cast", castHandler)
	http.HandleFunc("/screenshot.jpg", screenshotHandler)
	http.HandleFunc("/events", eventsHandler)
	fmt.Printf("Recorder running at http://localhost:8080/cast\n")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Printf("Player running at http://%s:8080\n", ipNet.IP.String())
			}
		}
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Index serves the static HTML content
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

// Cast handler to store the uploaded screenshot in memory
func castHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content into memory directly
	mu.Lock()
	imageMemory, err = io.ReadAll(file)
	if err != nil {
		mu.Unlock()
		http.Error(w, "Failed to store file in memory", http.StatusInternalServerError)
		return
	}

	// Notify clients in a non-blocking manner
	for client := range clients {
		go func(c chan []byte) {
			c <- imageMemory
		}(client)
	}
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Image stored in memory!")
}

// Screenshot handler to serve the most recent image
func screenshotHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if imageMemory == nil {
		http.Error(w, "No image available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageMemory)
}

// Events handler for SSE
func eventsHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	mu.Lock()
	clients[messageChan] = struct{}{}
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, messageChan)
		mu.Unlock()
		close(messageChan)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-messageChan:
			w.Write([]byte("data: "))
			encoded := base64.StdEncoding.EncodeToString(msg)
			w.Write([]byte(encoded))
			w.Write([]byte("\n\n"))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
