package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Garment API is running!")
}

func main() {
	//! registering the router
	http.HandleFunc("/health", healthCheckHandler)

	//! start the server
	port := ":8080"
	fmt.Printf("Starting Garment API server on port %s...\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server Crashed: ", err)
	}
}
