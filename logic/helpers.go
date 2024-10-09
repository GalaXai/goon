package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
)

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
