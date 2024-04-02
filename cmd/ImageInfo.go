package main

import (
	"image"
	"image/color"
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
