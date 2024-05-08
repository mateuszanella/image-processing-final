package main

import (
	"image"
	"image/color"
	"math"
	"sort"
)

type RGB struct {
	R, G, B uint32
}

type ImageInfo struct {
	Width  int
	Height int
	Pixels [][]RGB
}

type KernelType int

const (
	Cross KernelType = iota
	Square
	Diamond
)

func NewImageInfo(img image.Image) *ImageInfo {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([][]RGB, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]RGB, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[y][x] = RGB{R: r, G: g, B: b}
		}
	}

	return &ImageInfo{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (imgInfo *ImageInfo) GenerateImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := imgInfo.Pixels[y][x].R >> 8
			g := imgInfo.Pixels[y][x].G >> 8
			b := imgInfo.Pixels[y][x].B >> 8
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

// Basic
func (imgInfo *ImageInfo) NewBinary(threshold uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			grayValue := 0.299*float64(imgInfo.Pixels[y][x].R)/255.0 +
				0.587*float64(imgInfo.Pixels[y][x].G)/255.0 +
				0.114*float64(imgInfo.Pixels[y][x].B)/255.0

			if grayValue > float64(threshold) {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewGrayscale() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := imgInfo.Pixels[y][x].R >> 8
			g := imgInfo.Pixels[y][x].G >> 8
			b := imgInfo.Pixels[y][x].B >> 8
			gray := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			img.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	return img
}

// Basic Operations
func (imgInfo *ImageInfo) AddValue(value uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := min((imgInfo.Pixels[y][x].R>>8)+uint32(value), 255)
			g := min((imgInfo.Pixels[y][x].G>>8)+uint32(value), 255)
			b := min((imgInfo.Pixels[y][x].B>>8)+uint32(value), 255)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) SubtractValue(value uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := subtractWithLimit(imgInfo.Pixels[y][x].R>>8, value)
			g := subtractWithLimit(imgInfo.Pixels[y][x].G>>8, value)
			b := subtractWithLimit(imgInfo.Pixels[y][x].B>>8, value)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) MultiplyValue(value uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := min((imgInfo.Pixels[y][x].R>>8)*uint32(value), 255)
			g := min((imgInfo.Pixels[y][x].G>>8)*uint32(value), 255)
			b := min((imgInfo.Pixels[y][x].B>>8)*uint32(value), 255)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) DivideValue(value uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			if value == 0 {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			r := imgInfo.Pixels[y][x].R >> 8
			g := imgInfo.Pixels[y][x].G >> 8
			b := imgInfo.Pixels[y][x].B >> 8

			r = min(r/uint32(value), 255)
			g = min(g/uint32(value), 255)
			b = min(b/uint32(value), 255)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

// Filer Operations
func (imgInfo *ImageInfo) NewNegative(img2 *ImageInfo) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := 255 - (img2.Pixels[y][x].R >> 8)
			g := 255 - (img2.Pixels[y][x].G >> 8)
			b := 255 - (img2.Pixels[y][x].B >> 8)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewHistogramEqualization() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	histogram := make([]int, 256)
	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r, g, b := imgInfo.Pixels[y][x].R, imgInfo.Pixels[y][x].G, imgInfo.Pixels[y][x].B
			gray := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			histogram[gray]++
		}
	}

	cumHistogram := make([]int, 256) // Cumulative histogram
	cumHistogram[0] = histogram[0]
	for i := 1; i < 256; i++ {
		cumHistogram[i] = cumHistogram[i-1] + histogram[i]
	}

	totalPixels := imgInfo.Width * imgInfo.Height
	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r, g, b := imgInfo.Pixels[y][x].R, imgInfo.Pixels[y][x].G, imgInfo.Pixels[y][x].B
			grayValue := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			pixelValue := float64(cumHistogram[grayValue]) / float64(totalPixels)
			img.SetGray(x, y, color.Gray{Y: uint8(pixelValue * 255)})
		}
	}

	return img
}

// Spatial Domain Filters
func (imgInfo *ImageInfo) NewMeanFilter(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := uint32(0)
			g := uint32(0)
			b := uint32(0)
			count := 0

			for i := -size; i <= size; i++ {
				for j := -size; j <= size; j++ {
					if y+i < 0 || y+i >= imgInfo.Height || x+j < 0 || x+j >= imgInfo.Width {
						continue
					}

					r += imgInfo.Pixels[y+i][x+j].R
					g += imgInfo.Pixels[y+i][x+j].G
					b += imgInfo.Pixels[y+i][x+j].B
					count++
				}
			}

			r /= uint32(count)
			g /= uint32(count)
			b /= uint32(count)
			img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewMedianFilter(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := make([]uint32, 0)
			g := make([]uint32, 0)
			b := make([]uint32, 0)

			for i := -size; i <= size; i++ {
				for j := -size; j <= size; j++ {
					if y+i < 0 || y+i >= imgInfo.Height || x+j < 0 || x+j >= imgInfo.Width {
						continue
					}

					r = append(r, imgInfo.Pixels[y+i][x+j].R)
					g = append(g, imgInfo.Pixels[y+i][x+j].G)
					b = append(b, imgInfo.Pixels[y+i][x+j].B)
				}
			}

			sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
			sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })
			sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

			median := len(r) / 2
			img.Set(x, y, color.RGBA{uint8(r[median] >> 8), uint8(g[median] >> 8), uint8(b[median] >> 8), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewGaussianFilter(size int) *image.RGBA {
	if size%2 == 0 {
		size++
	}

	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	kernel := GaussianKernel(size, 1.0)

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := float64(0)
			g := float64(0)
			b := float64(0)

			for i := range kernel {
				for j := range kernel[i] {
					if y+i-size/2 < 0 || y+i-size/2 >= imgInfo.Height || x+j-size/2 < 0 || x+j-size/2 >= imgInfo.Width {
						continue
					}

					r += float64(imgInfo.Pixels[y+i-size/2][x+j-size/2].R>>8) * kernel[i][j]
					g += float64(imgInfo.Pixels[y+i-size/2][x+j-size/2].G>>8) * kernel[i][j]
					b += float64(imgInfo.Pixels[y+i-size/2][x+j-size/2].B>>8) * kernel[i][j]
				}
			}
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

// Morphological Operations
func (imgInfo *ImageInfo) NewDilation(size int, kernelType KernelType) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	kernel := GenerateKernel(size, kernelType)

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := uint32(0)
			g := uint32(0)
			b := uint32(0)

			for i := -size / 2; i <= size/2; i++ {
				for j := -size / 2; j <= size/2; j++ {
					if y+i < 0 || y+i >= imgInfo.Height || x+j < 0 || x+j >= imgInfo.Width {
						continue
					}

					if kernel[i+size/2][j+size/2] == 1 {
						r = max(r, imgInfo.Pixels[y+i][x+j].R)
						g = max(g, imgInfo.Pixels[y+i][x+j].G)
						b = max(b, imgInfo.Pixels[y+i][x+j].B)
					}
				}
			}

			img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewErosion(size int, kernelType KernelType) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	kernel := GenerateKernel(size, kernelType)

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := uint32(255)
			g := uint32(255)
			b := uint32(255)

			for i := -size / 2; i <= size/2; i++ {
				for j := -size / 2; j <= size/2; j++ {
					if y+i < 0 || y+i >= imgInfo.Height || x+j < 0 || x+j >= imgInfo.Width {
						continue
					}

					if kernel[i+size/2][j+size/2] == 1 {
						r = min(r, uint32(imgInfo.Pixels[y+i][x+j].R))
						g = min(g, uint32(imgInfo.Pixels[y+i][x+j].G))
						b = min(b, uint32(imgInfo.Pixels[y+i][x+j].B))
					}
				}
			}

			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

func (imgInfo *ImageInfo) NewContour() *image.RGBA {
	eroded := imgInfo.NewErosion(3, Square)
	contour := imgInfo.SubtractImage(NewImageInfo(eroded))

	return contour
}

func (imgInfo *ImageInfo) SubtractImage(img2 *ImageInfo) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := subtractWithLimit(imgInfo.Pixels[y][x].R, uint8(img2.Pixels[y][x].R))
			g := subtractWithLimit(imgInfo.Pixels[y][x].G, uint8(img2.Pixels[y][x].G))
			b := subtractWithLimit(imgInfo.Pixels[y][x].B, uint8(img2.Pixels[y][x].B))
			img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	return img
}

// Logical Operations
func (imgInfo *ImageInfo) NewNot(img2 *ImageInfo) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	for y := 0; y < imgInfo.Height; y++ {
		for x := 0; x < imgInfo.Width; x++ {
			r := 255 - (img2.Pixels[y][x].R >> 8)
			g := 255 - (img2.Pixels[y][x].G >> 8)
			b := 255 - (img2.Pixels[y][x].B >> 8)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return img
}

// Edge Detection
func (imgInfo *ImageInfo) NewPrewittFilter() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, imgInfo.Width, imgInfo.Height))

	kernelX := [][]int{
		{-1, 0, 1},
		{-1, 0, 1},
		{-1, 0, 1},
	}
	kernelY := [][]int{
		{-1, -1, -1},
		{0, 0, 0},
		{1, 1, 1},
	}

	for y := 1; y < imgInfo.Height-1; y++ {
		for x := 1; x < imgInfo.Width-1; x++ {
			grayX := 0
			grayY := 0

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					rgb := imgInfo.Pixels[y+i][x+j]
					gray := uint8(0.2126*float64(rgb.R>>8) + 0.7152*float64(rgb.G>>8) + 0.0722*float64(rgb.B>>8))

					grayX += int(gray) * kernelX[i+1][j+1]
					grayY += int(gray) * kernelY[i+1][j+1]
				}
			}

			magnitude := math.Sqrt(float64(grayX*grayX + grayY*grayY))
			normalizedMagnitude := magnitude / math.Sqrt(2*255*255) * 255
			gray := uint8(normalizedMagnitude)

			img.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	return img
}

// Helper functions
func subtractWithLimit(value uint32, limit uint8) uint32 {
	if value < uint32(limit) {
		return 0
	}
	return value - uint32(limit)
}

func GaussianKernel(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}

	var sum float64
	radius := size / 2
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			exp := -(float64(x*x+y*y) / (2 * sigma * sigma))
			kernel[x+radius][y+radius] = (1 / (2 * math.Pi * sigma * sigma)) * math.Exp(exp)
			sum += kernel[x+radius][y+radius]
		}
	}

	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

func GenerateKernel(size int, kernelType KernelType) [][]uint8 {
	kernel := make([][]uint8, size)
	for i := range kernel {
		kernel[i] = make([]uint8, size)
	}

	mid := size / 2

	switch kernelType {
	case Cross:
		for i := 0; i < size; i++ {
			kernel[mid][i] = 1
			kernel[i][mid] = 1
		}
	case Square:
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				kernel[i][j] = 1
			}
		}
	case Diamond:
		for i := 0; i <= mid; i++ {
			for j := mid - i; j <= mid+i; j++ {
				kernel[i][j] = 1
				kernel[size-i-1][j] = 1
			}
		}
	}

	return kernel
}
