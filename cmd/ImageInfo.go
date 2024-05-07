package main

import (
	"image"
	"image/color"
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

// Helper functions
func subtractWithLimit(value uint32, limit uint8) uint32 {
	if value < uint32(limit) {
		return 0
	}
	return value - uint32(limit)
}
