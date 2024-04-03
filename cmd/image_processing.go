package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"os"
)

func (app *Config) openImage(filename string) (image.Image, error) {
	file, err := os.Open("./storage/" + filename + ".jpg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (app *Config) saveImage(img image.Image, filename string) error {
	out, err := os.Create("./storage/" + filename)
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

func (app *Config) CreateGrayscale(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	grayImg := imgInfo.NewGrayscale()

	err = app.saveImage(grayImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateBinary(filename string, threshold uint8) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	binaryImg := imgInfo.NewBinary(threshold)

	err = app.saveImage(binaryImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Basic process
func (app *Config) TestImageManipulation(filename string) error {
	// Ill keep this function like this as an example of all the steps needed to do the stuff
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
