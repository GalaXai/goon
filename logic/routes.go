package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"
	"strings"

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
		log.Printf("\033[34mReceived request to %s\033[0m", r.URL.Path)
		if err := f(w, r); err != nil {
			log.Printf("\033[31mError handling request: %v\033[0m", err)
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		} else {
			log.Printf("\033[32mRequest to %s handled successfully\033[0m", r.URL.Path)
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
	router.HandleFunc("/load-image", makeHTTPHandleFunc(s.handleLoadImage))

	log.Println("JSON API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLoadImage(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	// Remove the data URL prefix if present
	base64Data := req.ImageBase64
	if prefix := "data:image/png;base64,"; strings.HasPrefix(base64Data, prefix) {
		base64Data = base64Data[len(prefix):]
	}

	// Decode base64 image
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("error decoding base64 image: %v", err)
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("error decoding image: %v", err)
	}

	// Convert image to matrix
	matrix := getImageMatrix(img)

	response := ImageResponse{
		Matrix: matrix,
		Shape:  img.Bounds().Max,
	}

	return WriteJSON(w, http.StatusOK, response)
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
