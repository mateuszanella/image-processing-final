package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image-processing/view/partials"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strconv"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"github.com/a-h/templ"
)

var kernelTypeMap = map[string]KernelType{
	"cross":   Cross,
	"square":  Square,
	"diamond": Diamond,
}

type BaseBody struct {
	Filename string `json:"filename"`
}

type BinaryBody struct {
	BaseBody
	Threshold string `json:"threshold"`
}

type BasicOperationsBody struct {
	BaseBody
	Value string `json:"value"`
}

type SpatialDomainBody struct {
	BaseBody
	Size string `json:"size"`
}

type OrderSpatialDomainBody struct {
	BaseBody
	Position string `json:"position"`
}

type MorphologicalOpeationsBody struct {
	BaseBody
	KernelType string `json:"kernelType"`
	Size       string `json:"size"`
}

// Image updates
func (app *Config) HandleStaticFiles() http.Handler {
	fs := http.FileServer(http.Dir("./static"))
	return http.StripPrefix("/static/", fs)
}

func (app *Config) HandleUploadImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to read image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, format, err := image.Decode(file)
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
			return
		}

		out, err := os.Create("./storage/uploaded." + format)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		if format == "tiff" || format == "tif" {
			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
				return
			}

			img, _, err := image.Decode(file)
			if err != nil {
				http.Error(w, "Failed to decode image", http.StatusInternalServerError)
				return
			}

			out2, err := os.Create("./storage/output.jpeg")
			if err != nil {
				http.Error(w, "Failed to open file", http.StatusInternalServerError)
				return
			}
			defer out2.Close()

			var opt jpeg.Options
			opt.Quality = 100

			err = jpeg.Encode(out2, img, &opt)
			if err != nil {
				http.Error(w, "Failed to encode image", http.StatusInternalServerError)
				return
			}
		} else {
			out2, err := os.Create("./storage/output." + format)
			if err != nil {
				http.Error(w, "Failed to open file", http.StatusInternalServerError)
				return
			}
			defer out2.Close()

			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(out2, file)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
		}

		templ.Handler(partials.ImageDisplay()).ServeHTTP(w, r)
	})
}

// On this upload, 5 files must be created:
// - The first uploaded image (and it's display (display-image1.jpg)) (image1.jpg)
// - The second uploaded image (and it's display (display-image2.jpg)) (image2.jpg)
// - The output image, for now, will be the last uploaded image (combination-output.jpg) just to make sure
func (app *Config) HandleCombinationUploadImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to read image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Wich of the images will be the upload
		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Operation parameter is required", http.StatusBadRequest)
			return
		}

		if filename != "image1" && filename != "image2" {
			http.Error(w, "Invalid operation parameter", http.StatusBadRequest)
			return
		}

		_, format, err := image.Decode(file)
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
			return
		}

		// Create the base image
		out, err := os.Create("./storage/" + filename + "." + format)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		// Create the display image
		if format == "tiff" || format == "tif" {
			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
				return
			}

			img, _, err := image.Decode(file)
			if err != nil {
				http.Error(w, "Failed to decode image", http.StatusInternalServerError)
				return
			}

			out2, err := os.Create("./storage/display-" + filename + ".jpeg")
			if err != nil {
				http.Error(w, "Failed to open file", http.StatusInternalServerError)
				return
			}
			defer out2.Close()

			var opt jpeg.Options
			opt.Quality = 100

			err = jpeg.Encode(out2, img, &opt)
			if err != nil {
				http.Error(w, "Failed to encode image", http.StatusInternalServerError)
				return
			}
		} else {
			out2, err := os.Create("./storage/display-" + filename + "." + format)
			if err != nil {
				http.Error(w, "Failed to open file", http.StatusInternalServerError)
				return
			}
			defer out2.Close()

			_, err = file.Seek(0, 0)
			if err != nil {
				http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(out2, file)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
		}

		templ.Handler(partials.ImageDisplay()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleGetImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filetype := r.URL.Query().Get("filetype")
		if filetype == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		filePath := "./storage/output." + filetype

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			filePath = "./static/uploaded.jpg"
			_, err = os.Stat(filePath)
			if os.IsNotExist(err) {
				http.Error(w, "Image not found", http.StatusNotFound)
				return
			}
		}

		if err != nil {
			http.Error(w, "Failed to open image", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, filePath)
	})
}

func (app *Config) HandleGetCombinationImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filetype := r.URL.Query().Get("filetype")
		if filetype == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		if filename != "image1" && filename != "image2" {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		filePath := "./storage/" + filename + "." + filetype

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to open image", http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, filePath)
	})
}

func (app *Config) HandleGetImageByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		file, err := os.Open("./storage/" + id + ".jpg")
		if err != nil {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		http.ServeFile(w, r, "./storage/"+id+".jpg")
	})
}

func (app *Config) HandleTestImageManipulation() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := app.TestImageManipulation("")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to manipulate image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Test ran sucessfully, check storage folder for output image")
	})
}

// Image filters

// Basic
func (app *Config) HandleCreateGrayscale() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateGrayscale(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create grayscale: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleCreateBinary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BinaryBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		threshold := 128
		if body.Threshold != "" {
			threshold, err = strconv.Atoi(body.Threshold)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse threshold: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateBinary(filename, uint8(threshold))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create binary: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Basic Operations
func (app *Config) HandleAddValue() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data BasicOperationsBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 10
		filename := ""
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
		}
		if data.Filename != "" {
			filename = data.Filename
		}

		if value < 0 || value > 255 {
			http.Error(w, "Value must be between 0 and 255", http.StatusBadRequest)
			return
		}

		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		err = app.AddPixels(filename, uint8(value))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add values: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleSubtractValue() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data BasicOperationsBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 10
		filename := ""
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
		}
		if data.Filename != "" {
			filename = data.Filename
		}

		if value < 0 || value > 255 {
			http.Error(w, "Value must be between 0 and 255", http.StatusBadRequest)
			return
		}

		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		err = app.SubtractPixels(filename, uint8(value))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to subtract values: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleMultiplyValue() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data BasicOperationsBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 10
		filename := ""
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
		}
		if data.Filename != "" {
			filename = data.Filename
		}

		if value < 0 || value > 255 {
			http.Error(w, "Value must be between 0 and 255", http.StatusBadRequest)
			return
		}

		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		err = app.MultiplyPixels(filename, uint8(value))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to multiply values: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleDivideValue() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data BasicOperationsBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 10
		filename := ""
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
		}
		if data.Filename != "" {
			filename = data.Filename
		}

		if value < 0 || value > 255 {
			http.Error(w, "Value must be between 0 and 255", http.StatusBadRequest)
			return
		}

		if filename == "" {
			http.Error(w, "Filename not provided", http.StatusBadRequest)
			return
		}

		err = app.DividePixels(filename, uint8(value))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to divide values: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Logical Operations
func (app *Config) HandleNotOpertion() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.NotOpertion(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply not operation: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Basic Filters
func (app *Config) HandleCreateNegativeFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateNegative(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create negative filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleCreateHisogramEqualization() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateHistogramEqualization(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create histogram equalization: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Spatial Domain Filters
func (app *Config) HandleMeanFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename
		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateMeanFilter(filename, size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply mean filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleMedianFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename
		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateMedianFilter(filename, size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply median filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleGaussianFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename
		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateGaussianFilter(filename, size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply gaussian filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleMinimumFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}
		filename := body.Filename
		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}
		err = app.CreateMinimumFilter(filename, size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply minimum filter: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleMaximumFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}
		filename := body.Filename
		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}
		err = app.CreateMaximumFilter(filename, size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply maximum filter: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleOrderFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body OrderSpatialDomainBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		filename := body.Filename
		position := 0
		if body.Position != "" {
			position, err = strconv.Atoi(body.Position)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse position: %v", err), http.StatusInternalServerError)
				return
			}

			if position < 0 || position > 8 {
				http.Error(w, "Position must be between 0 and 8", http.StatusBadRequest)
				return
			}
		}

		err = app.CreateOrderFilter(filename, position)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply order filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleConservativeSmoothingFilter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		filename := body.Filename

		err = app.CreateConservativeSmoothingFilter(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply conservative smoothing filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Morphological Operations
func (app *Config) HandleDilation() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body MorphologicalOpeationsBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		kernelType, err := getKernelTypeFromString(body.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateDilation(filename, size, kernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply dilation: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleErosion() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body MorphologicalOpeationsBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		kernelType, err := getKernelTypeFromString(body.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateErosion(filename, size, kernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply erosion: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleOpening() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body MorphologicalOpeationsBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		kernelType, err := getKernelTypeFromString(body.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateOpening(filename, size, kernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply opening: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleClosing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body MorphologicalOpeationsBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		kernelType, err := getKernelTypeFromString(body.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if body.Size != "" {
			size, err = strconv.Atoi(body.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse size: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateClosing(filename, size, kernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to apply closing: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleContour() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateContour(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create contour: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// Edge Detection
func (app *Config) HandlePrewittEdgeDetection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreatePrewittEdgeDetection(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create prewitt edge detection: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleSobelEdgeDetection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateSobelEdgeDetection(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create sobel edge detection: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleLaplacianEdgeDetection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateLaplacianEdgeDetection(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create sobel edge detection: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

// *-**-* Bonus *-**-*

// Image Adjustments
func (app *Config) HandleFlipLR() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateFlipLR(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to flip image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image flipped successfully, check storage folder for output image")
	})
}

func (app *Config) HandleFlipUD() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateFlipUD(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to flip image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image flipped successfully, check storage folder for output image")
	})
}

func (app *Config) HandleRotate90() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateRotate90(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to rotate image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image rotated successfully, check storage folder for output image")
	})
}

func (app *Config) HandleRotate270() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body BaseBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
			return
		}

		filename := body.Filename

		err = app.CreateRotate270(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to rotate image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image rotated successfully, check storage folder for output image")
	})
}

// Helpers
func getKernelTypeFromString(s string) (KernelType, error) {
	kernelType, ok := kernelTypeMap[s]
	if !ok {
		return 0, fmt.Errorf("invalid kernel type: %s", s)
	}
	return kernelType, nil
}

// ********** //
// Components
func (app *Config) HandleDisplayComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := partials.ImageDisplay()

		templ.Handler(c).ServeHTTP(w, r)
	})
}

func (app *Config) HandleDropzoneComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.Dropzone()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleCombinationDropzoneComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.CombinationDropzone()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleCombinationFiltersComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.CombinationFilters()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleFiltersComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.Filters()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleAdjustmentsComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.Adjustments()).ServeHTTP(w, r)
	})
}

func (app *Config) HandleAddImagesComponent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(partials.AddImagesForm()).ServeHTTP(w, r)
	})
}
