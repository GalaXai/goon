package main

import (
	"image"
)

type Matrix2D [][]uint8
type Matrix3D [][][]uint8
type RuneMatrix2D [][]rune
type RuneMatrix3D [][][]rune

type ImageRequest struct {
	ImageBase64 string `json:"image_base64"`
}

type ImageResponse struct {
	Matrix Matrix3D    `json:"matrix"`
	Shape  image.Point `json:"shape"`
}
