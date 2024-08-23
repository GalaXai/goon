package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
)

func getImageMatrix(img image.Image) [][][]uint8 {
	shape := img.Bounds().Max
	X, Y := shape.X, shape.Y

	/*
		Matrix should look like
		(3, 512, 512) where:
		3 is number of dims
		3 for RGB
		1 for black white
		512,512 is a 2D matrix that has value
	*/

	// Creates RGB part should work same with black white version for now, this could use optimization.
	matrix := make([][][]uint8, Y)
	// iter over [:]
	for i := range matrix {
		// Creates X, Y pixel values
		matrix[i] = make([][]uint8, X)
		for j := range matrix {
			R, G, B, _ := img.At(i, j).RGBA()
			// bitwise operation to normalize image values
			matrix[i][j] = []uint8{uint8(R >> 8), uint8(G >> 8), uint8(B >> 8)}
		}
	}
	return matrix
}

func downSample(matrix [][][]uint8, kernelSize int) [][][]uint8 {
	X, Y, dim := len(matrix), len(matrix[0]), len(matrix[0][0])

	newY := Y / kernelSize
	newX := X / kernelSize
	kernelArea := kernelSize * kernelSize

	downSampledMatrix := make([][][]uint8, newY)
	for i := range downSampledMatrix {
		downSampledMatrix[i] = make([][]uint8, newX)
		for j := range downSampledMatrix[i] {
			sum := make([]uint64, dim)
			// Gets average from the kernel window
			for ky := 0; ky < kernelSize; ky++ {
				y := i*kernelSize + ky
				for kx := 0; kx < kernelSize; kx++ {
					x := j*kernelSize + kx
					for d := 0; d < dim; d++ {
						sum[d] += uint64(matrix[y][x][d])
					}
				}
			}
			avg := make([]uint8, dim)
			for d := 0; d < dim; d++ {
				avg[d] = uint8(sum[d] / uint64(kernelArea))
			}
			downSampledMatrix[i][j] = avg
		}
	}
	return downSampledMatrix
}

func exportImage(matrix [][][]uint8, fileName string) {
	Y, X := len(matrix), len(matrix[1])
	imgCanvas := image.NewNRGBA(image.Rect(0, 0, X, Y))

	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			r := matrix[y][x][0]
			g := matrix[y][x][1]
			b := matrix[y][x][2]
			imgCanvas.Set(y, x, color.RGBA{r, g, b, 255}) // Assuming full opacity
		}
	}

	// Create a file to save the image
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Encode and save the image as PNG
	err = png.Encode(f, imgCanvas)
	if err != nil {
		panic(err)
	}
}

func desaturateInplace(matrix [][][]uint8) {

	Y, X := len(matrix), len(matrix[1])
	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			r := matrix[y][x][0]
			g := matrix[y][x][1]
			b := matrix[y][x][2]
			gray := float32(0.299)*float32(r) + float32(0.587)*float32(g) + float32(0.114)*float32(b)
			for d := 0; d < 3; d++ {
				matrix[y][x][d] = uint8(gray)
			}
		}
	}
}

func asciiIamge(matrix [][][]uint8) [][][]rune {
	Y, X := len(matrix), len(matrix[1])
	asciiTable := []rune{' ', ',', ';', 'c', 'o', 'P', 'O', '?', '@', '▓'}

	asciiMatrix := make([][][]rune, Y)
	for i := range asciiMatrix {
		asciiMatrix[i] = make([][]rune, X)
		for j := range asciiMatrix {
			luminance := uint8(int(math.Floor(float64(matrix[j][i][0]) / 255.0 * 9.0)))
			asciiMatrix[i][j] = make([]rune, 1)
			asciiMatrix[i][j][0] = asciiTable[luminance]
		}
	}
	return asciiMatrix
}

func printAsciiArt(asciiMatrix [][][]rune) {
	for _, row := range asciiMatrix {
		for _, cell := range row {
			fmt.Print(string(cell))
		}
		fmt.Println() // New line after each row
	}
}

// Uncomment for testing

func main() {
	//open the image
	file, err := os.Open("static/test.png")
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image :", err)
		return
	}
	matrix := getImageMatrix(img)
	dMatrix := downSample(matrix, 8)
	desaturateInplace(dMatrix)
	ascii := asciiIamge(dMatrix)
	printAsciiArt(ascii)
	exportImage(dMatrix, "static/downsapled.png")
}
