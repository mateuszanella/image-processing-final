package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/bmp"
)

func (app *Config) openImage(filename string) (image.Image, string, error) {
	file, err := os.Open("./storage/" + filename)
	if err != nil {
		fmt.Println("An error occurred while opening the image: ", err)
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)

	if err != nil {
		fmt.Println("An error occurred while decoding the image: ", err)
		return nil, "", err
	}

	return img, format, nil
}

func (app *Config) saveImage(img image.Image, format string) error {
	if format == "tiff" || format == "tif" {
		format = "jpeg"
	}

	out, err := os.Create("./storage/output." + format)
	if err != nil {
		fmt.Println("An error occurred while saving the image: ", err)
		return err
	}
	defer out.Close()

	switch format {
	case "jpeg":
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, img, &opt)
	case "png":
		err = png.Encode(out, img)
	case "bmp":
		err = bmp.Encode(out, img)
	default:
		err = fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		fmt.Println("An error occurred while encoding the image: ", err)
		return err
	}

	return nil
}

// fuck this
func (app *Config) saveOperationImage(img image.Image) error {
	out, err := os.Create("./storage/combination-output.jpeg")
	if err != nil {
		fmt.Println("An error occurred while saving the image: ", err)
		return err
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 100
	err = jpeg.Encode(out, img, &opt)

	if err != nil {
		fmt.Println("An error occurred while encoding the image: ", err)
		return err
	}

	return nil
}

// Basic filters
func (app *Config) CreateGrayscale(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	grayImg := imgInfo.NewGrayscale()

	err = app.saveImage(grayImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateBinary(filename string, threshold uint8) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	binaryImg := imgInfo.NewBinary(threshold)

	err = app.saveImage(binaryImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Basic Operations
func (app *Config) AddPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.AddValue(value)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) SubtractPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.SubtractValue(value)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) MultiplyPixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.MultiplyValue(value)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) DividePixels(filename string, value uint8) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.DivideValue(value)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Logical Operations
func (app *Config) NotOpertion(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := imgInfo.NewNot(binaryImgInfo)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) AndOperation(image1, image2 string) error {
	img1, _, err := app.openImage(image1)
	if err != nil {
		return err
	}

	img2, _, err := app.openImage(image2)
	if err != nil {
		return err
	}

	imgInfo1 := NewImageInfo(img1)
	binaryImg1 := imgInfo1.NewBinary(128)
	binaryImg1Info := NewImageInfo(binaryImg1)

	imgInfo2 := NewImageInfo(img2)
	binaryImg2 := imgInfo2.NewBinary(128)
	binaryImg2Info := NewImageInfo(binaryImg2)

	processedImg := binaryImg1Info.NewAnd(binaryImg2Info)

	err = app.saveOperationImage(processedImg)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) OrOperation(image1, image2 string) error {
	img1, _, err := app.openImage(image1)
	if err != nil {
		return err
	}

	img2, _, err := app.openImage(image2)
	if err != nil {
		return err
	}

	imgInfo1 := NewImageInfo(img1)
	binaryImg1 := imgInfo1.NewBinary(128)
	binaryImg1Info := NewImageInfo(binaryImg1)

	imgInfo2 := NewImageInfo(img2)
	binaryImg2 := imgInfo2.NewBinary(128)
	binaryImg2Info := NewImageInfo(binaryImg2)

	processedImg := binaryImg1Info.NewOr(binaryImg2Info)

	err = app.saveOperationImage(processedImg)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) XorOperation(image1, image2 string) error {
	img1, _, err := app.openImage(image1)
	if err != nil {
		return err
	}

	img2, _, err := app.openImage(image2)
	if err != nil {
		return err
	}

	imgInfo1 := NewImageInfo(img1)
	binaryImg1 := imgInfo1.NewBinary(128)
	binaryImg1Info := NewImageInfo(binaryImg1)

	imgInfo2 := NewImageInfo(img2)
	binaryImg2 := imgInfo2.NewBinary(128)
	binaryImg2Info := NewImageInfo(binaryImg2)

	processedImg := binaryImg1Info.NewXor(binaryImg2Info)

	err = app.saveOperationImage(processedImg)
	if err != nil {
		return err
	}

	return nil
}

// Filters
func (app *Config) CreateNegative(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	negativeImg := imgInfo.NewNegative(imgInfo)

	err = app.saveImage(negativeImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateHistogramEqualization(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	grayscaleImage := imgInfo.NewGrayscale()

	grayscaleImageInfo := NewImageInfo(grayscaleImage)
	processedImg := grayscaleImageInfo.NewHistogramEqualization()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Spatial Domain Filters
func (app *Config) CreateMeanFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMeanFilter(size)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateMedianFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMedianFilter(size)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateGaussianFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewGaussianFilter(size)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateMinimumFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMinFilter(size)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateMaximumFilter(filename string, size int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewMaxFilter(size)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateOrderFilter(filename string, position int) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewOrderFilter(position)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateConservativeSmoothingFilter(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}
	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewConservativeSmoothingFilter()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Morphological Operations
func (app *Config) CreateDilation(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewDilation(size, kernelType)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateErosion(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewErosion(size, kernelType)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateOpening(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	erosionImg := binaryImgInfo.NewErosion(size, kernelType)

	erosionImgInfo := NewImageInfo(erosionImg)
	processedImg := erosionImgInfo.NewDilation(size, kernelType)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateClosing(filename string, size int, kernelType KernelType) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	dilationImg := binaryImgInfo.NewDilation(size, kernelType)

	dilationImgInfo := NewImageInfo(dilationImg)
	processedImg := dilationImgInfo.NewErosion(size, kernelType)

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateContour(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)
	binaryImg := imgInfo.NewBinary(128)

	binaryImgInfo := NewImageInfo(binaryImg)
	processedImg := binaryImgInfo.NewContour()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Edge Detection
func (app *Config) CreatePrewittEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewPrewittFilter()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateSobelEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewSobelFilter()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateLaplacianEdgeDetection(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	imgInfo := NewImageInfo(img)

	processedImg := imgInfo.NewLaplacianFilter()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

// Image Based Operations

func (app *Config) CreateImageAddition(image1, image2 string) error {
	img1, _, err := app.openImage(image1)
	if err != nil {
		return err
	}

	img2, _, err := app.openImage(image2)
	if err != nil {
		return err
	}

	imgInfo1 := NewImageInfo(img1)
	imgInfo2 := NewImageInfo(img2)

	processedImg := imgInfo1.AddImages(imgInfo2)

	err = app.saveOperationImage(processedImg)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateImageSubtraction(image1, image2 string) error {
	img1, _, err := app.openImage(image1)
	if err != nil {
		return err
	}

	img2, _, err := app.openImage(image2)
	if err != nil {
		return err
	}

	imgInfo1 := NewImageInfo(img1)
	imgInfo2 := NewImageInfo(img2)

	processedImg := imgInfo1.SubtractImages(imgInfo2)

	err = app.saveOperationImage(processedImg)
	if err != nil {
		return err
	}

	return nil
}

// *-**-* Bonus *-**-*

// Image Adjustments
func (app *Config) CreateFlipLR(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	processedImg := NewImageInfo(img).FlipLR()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateFlipUD(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	processedImg := NewImageInfo(img).FlipUD()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateRotate90(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	processedImg := NewImageInfo(img).Rotate90()

	err = app.saveImage(processedImg, format)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) CreateRotate270(filename string) error {
	if filename == "" {
		filename = "uploaded.jpg"
	}

	img, format, err := app.openImage(filename)
	if err != nil {
		return err
	}

	processedImg := NewImageInfo(img).Rotate270()

	err = app.saveImage(processedImg, format)
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
