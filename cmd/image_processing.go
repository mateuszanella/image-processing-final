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

// Basic filters
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

// Basic Operations
func (app *Config) AddPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.AddValue(value)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) SubtractPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.SubtractValue(value)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) MultiplyPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.MultiplyValue(value)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) DividePixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.DivideValue(value)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Logical Operations
func (app *Config) NotOpertion(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := imgInfo.NewNot(binaryImgInfo)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Filters
func (app *Config) CreateNegative(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	negativeImg := imgInfo.NewNegative(imgInfo)

	err = app.saveImage(negativeImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateHistogramEqualization(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewHistogramEqualization()

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Spatial Domain Filters
func (app *Config) CreateMeanFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMeanFilter(size)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateMedianFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMedianFilter(size)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateGaussianFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewGaussianFilter(size)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Morphological Operations
func (app *Config) CreateDilation(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewDilation(size, kernelType)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateErosion(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewErosion(size, kernelType)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateOpening(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	erosionImg := binaryImgInfo.NewErosion(size, kernelType)

	erosionImgInfo := NewImageInfo(erosionImg)
	processedImg := erosionImgInfo.NewDilation(size, kernelType)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateClosing(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	dilationImg := binaryImgInfo.NewDilation(size, kernelType)

	dilationImgInfo := NewImageInfo(dilationImg)
	processedImg := dilationImgInfo.NewErosion(size, kernelType)

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateContour(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewContour()

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Edge Detection
func (app *Config) CreatePrewittEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewPrewittFilter()

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateSobelEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewSobelFilter()

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateLaplacianEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded"
	}

	img, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewLaplacianFilter()

	err = app.saveImage(processedImg, "output.jpg")
	if err != nil {
		return err
	}

	return nil
}

// Dummy
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
