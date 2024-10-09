package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
		// Blue text
		log.Printf("\033[34mReceived request to %s\033[0m", r.URL.Path)
		startTime := time.Now()
		if err := f(w, r); err != nil {
			// Red text
			log.Printf("\033[31mError handling request: %v\033[0m", err)
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		} else {
			duration := time.Since(startTime)
			// Green text for success, yellow text for duration
			log.Printf("\033[32mRequest to %s handled successfully\033[0m \033[33m(took %v)\033[0m", r.URL.Path, duration)
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

	// Generate a unique key for the image
	cacheKey := generateCacheKey(req.Base64Image)

	// Try to get the cached results
	cache := getImageCache(cacheKey)

	var originalMatrix, desaturatedMatrix, downsampledMatrix, gaussiansDiff Matrix3D
	var err error
	if cache == nil {
		originalMatrix, err = base64ToMatrix(req.Base64Image)
		if err != nil {
			return fmt.Errorf("error converting base64 to matrix: %v", err)
		}

		// Desaturate
		desaturatedMatrix = desaturate(originalMatrix)

		// Downsample
		downsampledMatrix, err = downSample(desaturatedMatrix, 1)
		if err != nil {
			return fmt.Errorf("error downsampling image: %v", err)
		}

		// Difference of Gaussians
		gaussiansDiff, err = differenceOfGaussians(downsampledMatrix)
		if err != nil {
			return fmt.Errorf("error applying difference of Gaussians: %v", err)
		}
		// Store the results in cache
		cache = &ImageCache{
			OriginalMatrix:    originalMatrix,
			DesaturatedMatrix: desaturatedMatrix,
			DownsampledMatrix: downsampledMatrix,
			GaussiansDiff:     gaussiansDiff,
			LastUsed:          time.Now(),
		}
		setImageCache(cacheKey, cache)
	} else {
		// Use cached results
		originalMatrix = cache.OriginalMatrix
		desaturatedMatrix = cache.DesaturatedMatrix
		downsampledMatrix = cache.DownsampledMatrix
		gaussiansDiff = cache.GaussiansDiff
		cache.LastUsed = time.Now()
	}
	// Sobel fiter -> __, edges(gradient at dataPoint)
	gradientThreshold := req.GradientThreshold
	if gradientThreshold <= 0 {
		gradientThreshold = 80
	}
	sobelMatrix, gradientMatrix, err := sobelFilter(gaussiansDiff, gradientThreshold)
	if err != nil {
		return fmt.Errorf("error in sobelFilter: %v", err)
	}

	// Create the Base64ImageResponse
	Base64ImageResponse := Base64ImageResponse{
		OriginalImage:      matrixToBase64(originalMatrix),
		DesaturatedImage:   matrixToBase64(desaturatedMatrix),
		DownsampledImage:   matrixToBase64(downsampledMatrix),
		GaussiansDiffImage: matrixToBase64(gaussiansDiff),
		SobelImage:         matrixToBase64(sobelMatrix),
		GradientImage:      matrixToBase64(gradientMatrix),
	}

	return WriteJSON(w, http.StatusOK, Base64ImageResponse)
}
