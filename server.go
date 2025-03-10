package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	service *Service
}

func NewServer() *Server {
	return &Server{
		service: NewService(),
	}
}

func (s *Server) incrementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	counterID, ok := vars["counter_id"]

	if ok {
		s.service.Increment(counterID)
		fmt.Fprintf(w, "OK")
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (s *Server) counterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	counterID, ok := vars["counter_id"]

	if ok {
		count := s.service.GetCounter(counterID)
		fmt.Fprintf(w, "%d", count)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
