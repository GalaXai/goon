package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
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
	router.HandleFunc("/edge-detect-ascii", makeHTTPHandleFunc(s.handleEdgeDetectAscii))
	router.HandleFunc("/color-downsample", makeHTTPHandleFunc(s.handleColorDownsample))

	// Add CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("JSON API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, corsHandler(router))
}

func (s *APIServer) handleColorDownsample(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req Base64ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	originalMatrix, err := base64ToMatrix(req.Base64Image)
	if err != nil {
		return fmt.Errorf("error converting base64 to matrix: %v", err)
	}

	// Downsample
	downsampledMatrix, err := downSample(originalMatrix, 4)
	if err != nil {
		return fmt.Errorf("error downsampling image: %v", err)
	}
	log.Print(int(downsampledMatrix[0][0][0]), int(downsampledMatrix[0][0][1]), int(downsampledMatrix[0][0][2]))
	// We convert to int since JsonWrite dosn't support uint8
	intMatrix := matrixToInt(downsampledMatrix)
	return WriteJSON(w, http.StatusOK, intMatrix)
}

func (s *APIServer) handleEdgeDetectAscii(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req Base64EdgeDetectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("error decoding request body: %v", err)
	}

	var originalMatrix, desaturatedMatrix, downsampledMatrix, gaussiansDiff Matrix3D
	var err error
	// if cache == nil {
	originalMatrix, err = base64ToMatrix(req.Base64Image)
	if err != nil {
		return fmt.Errorf("error converting base64 to matrix: %v", err)
	}

	// Desaturate
	desaturatedMatrix = desaturate(originalMatrix)

	// Downsample
	downsampledMatrix, err = downSample(desaturatedMatrix, 2)
	if err != nil {
		return fmt.Errorf("error downsampling image: %v", err)
	}

	// Difference of Gaussians
	log.Print(req.Tau, req.Threshold)
	gaussiansDiff, err = differenceOfGaussians(downsampledMatrix, req.Tau, req.Threshold)
	if err != nil {
		return fmt.Errorf("error applying extended Difference Of Gaussians: %v", err)
	}

	gradientThreshold := req.GradientThreshold

	horizontalSobel, err := horizontalSobel(gaussiansDiff)
	if err != nil {
		return fmt.Errorf("error in sobelFilter: %v", err)
	}
	verticalSobel, err := verticalSobel(gaussiansDiff, gradientThreshold)
	if err != nil {
		return fmt.Errorf("error in sobelFilter: %v", err)
	}
	// Create the Base64ImageResponse
	Base64ImageResponse := Base64ImagesResponse{
		OriginalImage:      matrixToBase64(originalMatrix),
		DesaturatedImage:   matrixToBase64(desaturatedMatrix),
		DownsampledImage:   matrixToBase64(downsampledMatrix),
		GaussiansDiffImage: matrixToBase64(gaussiansDiff),
		HorizontalSobel:    matrixToBase64(horizontalSobel),
		VerticalSobel:      matrixToBase64(verticalSobel),
	}
	ascii_1 := asciiImage(verticalSobel, true, 0)
	// ascii_1 := asciiImage(downsampledMatrix, false, 0.1)
	ascii_2 := asciiImage(downsampledMatrix, false, 0.1)
	merged_ascii := mergeAsciiImages(ascii_1, ascii_2)

	runeMatrix3D := RuneMatrix3D{
		Data:  merged_ascii,
		Cols:  len(merged_ascii[0]),
		Rows:  len(merged_ascii),
		Depth: len(merged_ascii[0][0]),
	}

	combinedResponse := CombinedResponse{
		ImageResponse: Base64ImageResponse,
		AsciiArt:      runeMatrix3D,
	}
	return WriteJSON(w, http.StatusOK, combinedResponse)
}
