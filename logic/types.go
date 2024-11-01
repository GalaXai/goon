package main

import (
	"sync"
	"time"
)

type Matrix2D [][]uint8
type Matrix3D [][][]uint8

type RuneMatrix2D [][]rune
type RuneMatrix3D struct {
	Data  [][][]rune `json:"data"`
	Cols  int        `json:"cols"`
	Rows  int        `json:"rows"`
	Depth int        `json:"depth"`
}
type Base64ImageRequest struct {
	Base64Image string `json:"base64Image"`
}
type Base64EdgeDetectionRequest struct {
	Base64Image       string  `json:"base64Image"`
	GradientThreshold uint16  `json:"gradientThreshold"`
	Tau               float64 `json:"tau"`
	Threshold         float64 `json:"threshold"`
}
type CombinedResponse struct {
	ImageResponse Base64ImagesResponse `json:"imageResponse"`
	AsciiArt      RuneMatrix3D         `json:"asciiArt"`
}

var imageCaches = make(map[string]*ImageCache)
var cachesMutex sync.RWMutex

type ImageCache struct {
	OriginalMatrix    Matrix3D
	DesaturatedMatrix Matrix3D
	DownsampledMatrix Matrix3D
	GaussiansDiff     Matrix3D
	LastUsed          time.Time
	mutex             sync.RWMutex
}

type Base64ImageResponse struct {
	image string `json:"image"`
}

type ImagesResponse struct {
	OriginalMatrix      Matrix3D `json:"originalMatrix"`
	DesaturatedMatrix   Matrix3D `json:"desaturatedMatrix"`
	DownsampledMatrix   Matrix3D `json:"downsampledMatrix"`
	GaussiansDiffMatrix Matrix3D `json:"gaussiansDiffMatrix"`
	HorizontalSobel     Matrix3D `json:"horizontalSobel"`
	VerticalSobel       Matrix3D `json:"verticalSobel"`
}

type Base64ImagesResponse struct {
	OriginalImage      string `json:"originalImage"`
	DesaturatedImage   string `json:"desaturatedImage"`
	DownsampledImage   string `json:"downsampledImage"`
	GaussiansDiffImage string `json:"gaussiansDiffImage"`
	HorizontalSobel    string `json:"horizontalSobel"`
	VerticalSobel      string `json:"verticalSobel"`
}
