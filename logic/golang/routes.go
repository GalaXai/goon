package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle the error
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	// Our host
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/image", makeHTTPHandleFunc(s.handleImage))

	log.Println("JSON API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleImage(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetImage(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateImage(w, r) // no input?
	}
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetImage(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleCreateImage(w http.ResponseWriter, r *http.Request) error {
	return nil
}
