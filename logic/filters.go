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
		for j := range matrix[i] {
			R, G, B, _ := img.At(j, i).RGBA()
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
	fmt.Println(Y, X)
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
	Y, X, dim = len(downSampledMatrix), len(downSampledMatrix[0]), len(downSampledMatrix[0][0])
	fmt.Println(Y, X)
	return downSampledMatrix, nil
}

func exportImage(matrix interface{}, fileName string) {
	var Y, X int
	var imgCanvas *image.NRGBA

	switch m := matrix.(type) {
	case [][][]uint8:
		Y, X = len(m), len(m[0])
		imgCanvas = image.NewNRGBA(image.Rect(0, 0, X, Y))
		for i := 0; i < Y; i++ {
			for j := 0; j < X; j++ {
				r := m[i][j][0]
				g := m[i][j][1]
				b := m[i][j][2]
				imgCanvas.Set(j, i, color.RGBA{r, g, b, 255}) // Assuming full opacity
			}
		}
	case [][][]float64:
		Y, X = len(m), len(m[0])
		imgCanvas = image.NewNRGBA(image.Rect(0, 0, X, Y))
		for i := 0; i < Y; i++ {
			for j := 0; j < X; j++ {
				r := uint8(math.Min(math.Max(m[i][j][0]*255.0, 0), 255))
				g := uint8(math.Min(math.Max(m[i][j][1]*255.0, 0), 255))
				b := uint8(math.Min(math.Max(m[i][j][2]*255.0, 0), 255))
				imgCanvas.Set(j, i, color.RGBA{r, g, b, 255}) // Assuming full opacity
			}
		}
	default:
		panic("Unsupported matrix type")
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
			if len(matrix[i][j]) == 0 {
				continue
			}
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

func desaturate(matrix [][][]uint8) [][][]uint8 {
	Y, X := len(matrix), len(matrix[0])
	desaturatedMatrix := make([][][]uint8, Y)

	for i := 0; i < Y; i++ {
		desaturatedMatrix[i] = make([][]uint8, X)
		for j := 0; j < X; j++ {
			if len(matrix[i][j]) == 0 {
				desaturatedMatrix[i][j] = []uint8{}
				continue
			}
			r := matrix[i][j][0]
			g := matrix[i][j][1]
			b := matrix[i][j][2]
			gray := uint8(float32(0.299)*float32(r) + float32(0.587)*float32(g) + float32(0.114)*float32(b))
			desaturatedMatrix[i][j] = []uint8{gray, gray, gray}
		}
	}

	return desaturatedMatrix
}

// ▚

// ▞

// ▮
// ▬
// ▐
// ■
// ◢◣◤◥
// ⚊ ❙
func angleAsciiChar(angle uint8, MAGNITUDE_THRESHOLD float64) rune {
	f_angle := float64(angle) / 255
	// asciiValues := []rune{' ', '|', '\\', '/', '_'} // List of values corresponding to angles
	// asciiValues := []rune{' ', '▮', '▚', '▞', '▬'} // List of values corresponding to angles
	asciiValues := []rune{' ', '▮', '╲', '╱', '▬'} // List of values corresponding to angles

	switch {
	case f_angle == 0.5-MAGNITUDE_THRESHOLD || f_angle == 0.5+MAGNITUDE_THRESHOLD:
		return asciiValues[4] // Return '_' character
	case f_angle > 0.5+MAGNITUDE_THRESHOLD:
		return asciiValues[3] // Return '/' character
	case f_angle < 0.5-MAGNITUDE_THRESHOLD && f_angle > MAGNITUDE_THRESHOLD:
		return asciiValues[2] // Return '\' character
	case f_angle > 1-MAGNITUDE_THRESHOLD:
		return asciiValues[1] // Return '|' character
	default:
		return asciiValues[0] // Return ' ' character
	}
}

func getAsciiChar(value float64, angled bool, MAGNITUDE_THRESHOLD float64) rune {
	// asciiTable := []rune{' ', ',', ';', 'c', 'o', 'P', 'O', '?', '@', '▓'}
	asciiTable := []rune{' ', '░', '▒', '▓', '█'}
	if angled {
		return angleAsciiChar(uint8(value*255), MAGNITUDE_THRESHOLD)
	}
	luminance := uint8(int(math.Floor(value * float64(len(asciiTable)-1))))
	// fmt.Println(luminance, value, len(asciiTable))
	return asciiTable[luminance]
}

func asciiImage(matrix [][][]uint8, angled bool, MAGNITUDE_THRESHOLD float64) [][][]rune {
	Y, X := len(matrix), len(matrix[0])
	asciiMatrix := make([][][]rune, Y)
	for i := range asciiMatrix {
		asciiMatrix[i] = make([][]rune, X)
		for j := range asciiMatrix[i] {
			luminance := float64(matrix[i][j][1]) / 255.0 // TODO: For now we are getting 2nd value just for vertical sobel adjust later
			asciiMatrix[i][j] = []rune{getAsciiChar(luminance, angled, MAGNITUDE_THRESHOLD)}
		}
	}
	return asciiMatrix
}

func horizontalSobel(matrix [][][]uint8) ([][][]uint8, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in horizontalSobel:", r)
		}
	}()

	// Check if matrix is empty
	if len(matrix) == 0 || len(matrix[0]) == 0 || len(matrix[0][0]) == 0 {
		return nil, fmt.Errorf("input matrix is empty")
	}

	Y, X, dim := len(matrix), len(matrix[0]), len(matrix[0][0])
	result := make([][][]uint8, Y)
	for i := range result {
		result[i] = make([][]uint8, X)
		for j := range result[i] {
			result[i][j] = make([]uint8, 3) // Gx and Gy and 0
		}
	}

	for y := 0; y < Y; y++ {
		for x := 1; x < X-1; x++ {
			for d := 0; d < dim; d++ {
				lum1 := int(matrix[y][x-1][d])
				lum2 := int(matrix[y][x][d])
				lum3 := int(matrix[y][x+1][d])

				Gx := 3*lum1 + 0*lum2 + -3*lum3
				Gy := 3*lum1 + 0*lum2 + -3*lum3

				result[y][x][0] = uint8(Clamp(Gx, 0, 255))
				result[y][x][1] = uint8(Clamp(Gy, 0, 255))
			}
		}
	}

	return result, nil
}

func verticalSobel(horizontalGradient [][][]uint8, MAGNITUDE_THRESHOLD uint16) ([][][]uint8, error) {
	Y, X := len(horizontalGradient), len(horizontalGradient[0])
	result := make([][][]uint8, Y)
	for i := range result {
		result[i] = make([][]uint8, X)
		for j := range result[i] {
			result[i][j] = make([]uint8, 3) // magnitude and angle and 0
		}
	}

	for y := 1; y < Y-1; y++ {
		for x := 0; x < X; x++ {
			grad1 := horizontalGradient[y-1][x]
			grad2 := horizontalGradient[y][x]
			grad3 := horizontalGradient[y+1][x]

			Gx := 3*float64(grad1[0]) + 10*float64(grad2[0]) + 3*float64(grad3[0])
			Gy := 3*float64(grad1[1]) + 0*float64(grad2[1]) + -3*float64(grad3[1])

			magnitude := uint8(math.Min(255.0, math.Abs(Gx)+math.Abs(Gy)))
			var angle uint8
			// Calculate angle and bin @git
			if uint16(magnitude) > MAGNITUDE_THRESHOLD {
				// Convert angle to 0-255 range
				angle = uint8((math.Atan2(float64(Gy), float64(Gx)) + math.Pi) / (2 * math.Pi) * 255)
			} else {
				angle = 0 // Special value for non-edge pixels
			}
			result[y][x][0] = magnitude * 0 // TODO: REMOVE LAYER FOR NOW DOSN'T GIVE ANY INFO ON RETURN IMAGE
			result[y][x][1] = angle
		}
	}

	return result, nil
}

func sobelFilter(matrix [][][]uint8, MAGNITUDE_THRESHOLD uint16) ([][][]uint8, [][][]uint8, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic in sobelFilter:", r)
		}
	}()

	// Check if matrix is empty
	if len(matrix) == 0 || len(matrix[0]) == 0 || len(matrix[0][0]) == 0 {
		return nil, nil, fmt.Errorf("input matrix is empty")
	}

	Y, X, dim := len(matrix), len(matrix[0]), len(matrix[0][0])

	// Define Sobel kernels
	kernelX := [][]float64{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	kernelY := [][]float64{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	// Create result matrices
	sobelMatrix := make([][][]uint8, Y)
	gradientMatrix := make([][][]uint8, Y)
	for i := range sobelMatrix {
		sobelMatrix[i] = make([][]uint8, X)
		gradientMatrix[i] = make([][]uint8, X)
		for j := range sobelMatrix[i] {
			sobelMatrix[i][j] = make([]uint8, dim)
			gradientMatrix[i][j] = make([]uint8, dim)
		}
	}

	// Apply Sobel filter
	for y := 1; y < Y-1; y++ {
		for x := 1; x < X-1; x++ {
			for d := 0; d < dim; d++ {
				gx, gy := 0.0, 0.0
				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						if len(matrix[y+i][x+j]) == 0 {
							continue
						}
						pixel := float64(matrix[y+i][x+j][d])

						gx += pixel * kernelX[i+1][j+1]
						gy += pixel * kernelY[i+1][j+1]
					}
				}
				// magnitude := uint8(math.Abs(gx) + math.Abs(gy))
				magnitude := uint8(math.Min(255.0, math.Abs(gx)+math.Abs(gy)))
				sobelMatrix[y][x][d] = magnitude

				// Calculate angle and bin @git
				if uint16(magnitude) > MAGNITUDE_THRESHOLD {
					angle := math.Atan2(gy, gx)
					// Normalize angle to [0, 1) range
					binned := uint8(angle/math.Pi*0.5 + 0.5)
					gradientMatrix[y][x][d] = binned
				} else {
					gradientMatrix[y][x][d] = 128 // Special value for non-edge pixels
				}
			}
		}
	}

	return sobelMatrix, gradientMatrix, nil
}

func differenceOfGaussians(matrix [][][]uint8, tau float64, threshold float64) ([][][]uint8, error) {
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
				G_sigma := float64(blurred2[i][j][d])
				G_k_sigma := float64(blurred1[i][j][d])

				diff := G_sigma - tau*G_k_sigma
				// TODO : thest thresholding and 1+tau seems good to get insides of the pictures.
				// diff := (1+tau)*G_sigma - tau*G_k_sigma

				// if diff >= threshold {
				// 	diff = 255 // Set to maximum value if above threshold
				// } else {
				// 	diff = 0 // Set to minimum value if below threshold
				// }

				result[i][j][d] = uint8(Clamp(int(diff), 0, 255))
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
	}
	// normalize the kernel
	for i := range kernel {
		kernel[i] /= sum
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

func mergeAsciiImages(ascii_1, ascii_2 [][][]rune) [][][]rune {
	Y, X := len(ascii_1), len(ascii_1[0])
	merged := make([][][]rune, Y)
	for i := range merged {
		merged[i] = make([][]rune, X)
		for j := range merged[i] {
			if string(ascii_1[i][j][0]) == " " {
				merged[i][j] = ascii_2[i][j] // Use value from second ASCII image
			} else {
				merged[i][j] = ascii_1[i][j] // Default to first ASCII image
			}
		}
	}
	return merged
}
