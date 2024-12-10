package main

import (
	"fmt"
	"log"
	"net/http"
	"travel-planner/cmd"
)

func main() {
	// Start the HTTP server
	http.HandleFunc("/plan", cmd.HandlePlanRequest) // Define a handler for the /plan endpoint
	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
