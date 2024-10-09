package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
)

func isBase64Image(input string) bool {
	// Check if the input starts with the base64 image prefix
	if strings.HasPrefix(input, "data:image/") && strings.Contains(input, ";base64,") {
		return true
	}

	// If no prefix, try to decode and see if it's valid base64
	_, err := base64.StdEncoding.DecodeString(input)
	return err == nil
}

func base64ToMatrix(base64Data string) (Matrix3D, error) {
	// Remove the data URL prefix if present
	if prefix := "data:image/png;base64,"; strings.HasPrefix(base64Data, prefix) {
		base64Data = base64Data[len(prefix):]
	}

	// Decode base64 image
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 image: %v", err)
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	// Convert image to matrix
	return getImageMatrix(img), nil
}

func matrixToBase64(matrix Matrix3D) string {
	// Implementation depends on how you want to convert the matrix to an image
	// This is a placeholder function
	img := matrixToImage(matrix)
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func matrixToImage(matrix Matrix3D) image.Image {
	Y, X := len(matrix), len(matrix[0])
	imgCanvas := image.NewNRGBA(image.Rect(0, 0, X, Y))

	for i := 0; i < Y; i++ {
		for j := 0; j < X; j++ {
			r := matrix[i][j][0]
			g := matrix[i][j][1]
			b := matrix[i][j][2]
			imgCanvas.Set(j, i, color.RGBA{r, g, b, 255}) // Assuming full opacity
		}
	}

	return imgCanvas
}
