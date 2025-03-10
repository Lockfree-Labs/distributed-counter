package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for local development
	}

	server := NewServer()
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Handle("/counter/{counter_id}/increment", http.HandlerFunc(server.incrementHandler))
	router.Handle("/counter/{counter_id}", http.HandlerFunc(server.counterHandler))

	log.Println(fmt.Sprintf("Starting server on :%s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
