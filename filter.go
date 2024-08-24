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

func downSample(matrix [][][]uint8, kernelSize int) ([][][]uint8, error) {
	// Recover from panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in downSample:", r)
		}
	}()

	// Check if matrix is empty
	if len(matrix) == 0 || len(matrix[0]) == 0 || len(matrix[0][0]) == 0 {
		return nil, fmt.Errorf("input matrix is empty")
	}

	Y, X, dim := len(matrix), len(matrix[0]), len(matrix[0][0])

	// Check if kernelSize is valid
	if kernelSize <= 0 || kernelSize > Y || kernelSize > X {
		return nil, fmt.Errorf("invalid kernel size: %d", kernelSize)
	}

	// Adjust Y and X to be divisible by kernelSize
	Y = (Y / kernelSize) * kernelSize
	X = (X / kernelSize) * kernelSize

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
					// Check if indices are within bounds
					if y >= Y || x >= X {
						return nil, fmt.Errorf("index out of range: [%d][%d]", y, x)
					}
					if len(matrix[y][x]) == 0 {
						continue
					}
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
	return downSampledMatrix, nil
}

func exportImage(matrix [][][]uint8, fileName string) {
	Y, X := len(matrix), len(matrix[1])
	imgCanvas := image.NewNRGBA(image.Rect(0, 0, X, Y))

	for i := 0; i < Y; i++ {
		for j := 0; j < X; j++ {
			r := matrix[i][j][0]
			g := matrix[i][j][1]
			b := matrix[i][j][2]
			imgCanvas.Set(i, j, color.RGBA{r, g, b, 255}) // Assuming full opacity
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
	for i := 0; i < Y; i++ {
		for j := 0; j < X; j++ {
			r := matrix[i][j][0]
			g := matrix[i][j][1]
			b := matrix[i][j][2]
			gray := float32(0.299)*float32(r) + float32(0.587)*float32(g) + float32(0.114)*float32(b)
			for d := 0; d < 3; d++ {
				matrix[i][j][d] = uint8(gray)
			}
		}
	}
}

func asciiIamge(matrix [][][]uint8) [][][]rune {
	Y, X := len(matrix), len(matrix[0])
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

func differenceOfGaussians(matrix [][][]uint8) ([][][]uint8, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in differenceOfGaussians:", r)
		}
	}()

	// Check if matrix is empty
	if len(matrix) == 0 || len(matrix[0]) == 0 || len(matrix[0][0]) == 0 {
		return nil, fmt.Errorf("input matrix is empty")
	}

	Y, X, dim := len(matrix), len(matrix[0]), len(matrix[0][0])

	sigma1, sigma2 := 1.0, 1.6
	// kernel is list of values to multiplay a matrix by.
	kernel1, kernel2 := gaussianKernel(sigma1), gaussianKernel(sigma2)

	// matmul matrix over kernel values
	blurred1 := applyGaussianBlur(matrix, kernel1)
	blurred2 := applyGaussianBlur(matrix, kernel2)

	// blurred2 - blurred1
	result := make([][][]uint8, Y)
	for i := 0; i < Y; i++ {
		result[i] = make([][]uint8, X)
		for j := 0; j < X; j++ {
			result[i][j] = make([]uint8, dim)
			for d := 0; d < dim; d++ {
				diff := int(blurred2[i][j][d]) - int(blurred1[i][j][d])
				result[i][j][d] = uint8(Clamp(diff, 0, 255))
			}
		}
	}

	return result, nil

}

func gaussianKernel(sigma float64) []float64 {
	/*
		size of kernel -> Round sigma to neares int and multiply by 6
		also ensure its odd
	*/

	size := int(math.Ceil(6 * sigma))
	if size%2 == 0 {
		size++
	}
	kernel := make([]float64, size)
	sum := 0.0
	center := size / 2

	for i := 0; i < size; i++ {
		x := float64(i - center)
		/* kernel[i] = exp(-(x^2) / (2 * sigma^2)) / (sqrt(2 * pi) * sigma) */
		kernel[i] = math.Exp(-(x*x)/(2*sigma*sigma)) / (math.Sqrt(2*math.Pi) * sigma)
		sum += kernel[i]

		// normalize the kernel
		for i := range kernel {
			kernel[i] /= sum
		}

	}
	return kernel
}

func applyGaussianBlur(matrix [][][]uint8, kernel []float64) [][][]uint8 {
	Y, X, dim := len(matrix), len(matrix[0]), len(matrix[0][0])
	kernelSize := len(kernel)
	kernelRadius := kernelSize / 2

	temp := make([][][]uint8, Y)
	for i := range temp {
		temp[i] = make([][]uint8, X)
		for j := range temp[i] {
			temp[i][j] = make([]uint8, dim)
		}
	}
	// Apply horizontal blur ----

	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			for d := 0; d < dim; d++ {
				sum := 0.0       // float64
				weightSum := 0.0 // float64
				for i := 0; i < kernelSize; i++ {
					kx := x + i - kernelRadius
					if kx >= 0 && kx < X {
						if len(matrix[y][kx]) == 0 {
							continue
						}
						sum += float64(matrix[y][kx][d]) * kernel[i]
						weightSum += kernel[i]
					}
				}
				// Normalize the sum by dividing by the total weight
				sum /= weightSum
				temp[y][x][d] = uint8(Clamp(int(sum), 0, 255))
			}
		}
	}
	// Apply vertical blur \
	gaussianBlurMatrix := make([][][]uint8, Y)
	for i := range gaussianBlurMatrix {
		gaussianBlurMatrix[i] = make([][]uint8, X)
		for j := range gaussianBlurMatrix[i] {
			gaussianBlurMatrix[i][j] = make([]uint8, dim)
		}
	}

	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			for d := 0; d < dim; d++ {
				sum := 0.0
				weightSum := 0.0 // float64
				for i := 0; i < kernelSize; i++ {
					ky := y + i - kernelRadius
					if ky >= 0 && ky < Y {
						weightSum += kernel[i]
						sum += float64(temp[ky][x][d]) * kernel[i]
					}
				}
				sum /= weightSum
				gaussianBlurMatrix[y][x][d] = uint8(Clamp(int(sum), 0, 255))
			}
		}
	}

	return gaussianBlurMatrix
}

func Clamp(value int, min int, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value

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
	fmt.Println("Image shape Y:", len(matrix), "X:", len(matrix[0]))

	dMatrix, err := downSample(matrix, 8)
	if err != nil {
		fmt.Println("Error in downSample:", err)
		// Handle the error appropriately
		return
	}
	exportImage(dMatrix, "static/downsapled.png")
	fmt.Println("Image shape Y:", len(dMatrix), "X:", len(dMatrix[0]))

	desaturateInplace(dMatrix)
	exportImage(dMatrix, "static/desaturated.png")
	fmt.Println("Image shape Y:", len(dMatrix), "X:", len(dMatrix[0]))

	ascii := asciiIamge(dMatrix)
	fmt.Println("Image shape Y:", len(ascii), "X:", len(ascii[0]))
	printAsciiArt(ascii)

	gaussiansDiff, err := differenceOfGaussians(matrix)
	if err != nil {
		fmt.Println("Error in gaussiansDiff:", err)
	}
	exportImage(gaussiansDiff, "static/gaussDiff.png")
}
