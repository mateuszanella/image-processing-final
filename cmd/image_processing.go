package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"os"
)

func (app *Config) TestImageManipulation(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	//Open the image file

	file, err := os.Open("./storage/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	//Get image pixels

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([][]color.Color, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]color.Color, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = img.At(x, y)
		}
	}

	//Recreate the image

	processedImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := pixels[y][x]
			r, g, b, a := c.RGBA()
			processedImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	//Save the image

	out, err := os.Create("./storage/output.jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 100
	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		return err
	}

	return nil
}
