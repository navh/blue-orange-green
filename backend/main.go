// main.go
package main

import (
	"io"
	"log"
	"net/http"

	pb "buoyboy/proto"

	"google.golang.org/protobuf/proto"
)

func main() {
	http.HandleFunc("/buoy", handleBuoyReading)
	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBuoyReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}


	message := &pb.BuoyStatus{} // Update this to match your protobuf message type
	if err := proto.Unmarshal(body, message); err != nil {
		http.Error(w, "Error parsing protobuf message", http.StatusBadRequest)
		return
	}

	// Just print the message
	log.Printf("Received buoy reading: %+v", message)

	w.WriteHeader(http.StatusOK)
}
