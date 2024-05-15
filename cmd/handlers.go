package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image-processing/view/partials"
	_ "image/gif"
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

var formatMap = map[string]string{
	"jpeg": "jpg",
	"png":  "png",
	"bmp":  "bmp",
	"tiff": "tiff",
	"gif":  "gif",
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

		if val, ok := formatMap[format]; ok {
			format = val
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

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "Failed to reset file reader", http.StatusInternalServerError)
			return
		}

		out2, err := os.Create("./storage/output." + format)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer out2.Close()

		_, err = io.Copy(out2, file)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
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

		if val, ok := formatMap[filetype]; ok {
			filetype = val
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
		filename := r.FormValue("output")
		err := app.CreateGrayscale(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create grayscale: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created sucessfully, check storage folder for output image")
	})
}

func (app *Config) HandleCreateBinary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")

		var data struct {
			Threshold string `json:"threshold"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		threshold := 128
		if data.Threshold != "" {
			threshold, err = strconv.Atoi(data.Threshold)
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
		filename := r.FormValue("image")

		var data struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 10
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
		}

		if value < 0 || value > 255 {
			http.Error(w, "Value must be between 0 and 255", http.StatusBadRequest)
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
		filename := r.FormValue("image")

		var data struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 128
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
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
		filename := r.FormValue("image")

		var data struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 128
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
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
		filename := r.FormValue("image")

		var data struct {
			Value string `json:"value"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		value := 128
		if data.Value != "" {
			value, err = strconv.Atoi(data.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse value: %v", err), http.StatusInternalServerError)
				return
			}
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
		filename := r.FormValue("image")
		err := app.NotOpertion(filename)
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
		filename := r.FormValue("image")
		err := app.CreateNegative(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create negative filter: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleCreateHisogramEqualization() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateHistogramEqualization(filename)
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
		filename := r.FormValue("image")

		var data struct {
			Size string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")

		var data struct {
			Size string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")

		var data struct {
			Size string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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

// Morphological Operations
func (app *Config) HandleDilation() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")

		var data struct {
			KernelType string `json:"kernelType"`
			Size       string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		kernelType, err := getKernelTypeFromString(data.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")

		var data struct {
			KernelType string `json:"kernelType"`
			Size       string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		kernelType, err := getKernelTypeFromString(data.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")

		var data struct {
			KernelType string `json:"kernelType"`
			Size       string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		kernelType, err := getKernelTypeFromString(data.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")

		var data struct {
			KernelType string `json:"kernelType"`
			Size       string `json:"size"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusInternalServerError)
			return
		}

		kernelType, err := getKernelTypeFromString(data.KernelType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse kernel type: %v", err), http.StatusInternalServerError)
			return
		}

		size := 3
		if data.Size != "" {
			size, err = strconv.Atoi(data.Size)
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
		filename := r.FormValue("image")
		err := app.CreateContour(filename)
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
		filename := r.FormValue("image")
		err := app.CreatePrewittEdgeDetection(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create prewitt edge detection: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleSobelEdgeDetection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateSobelEdgeDetection(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create sobel edge detection: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created successfully, check storage folder for output image")
	})
}

func (app *Config) HandleLaplacianEdgeDetection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateLaplacianEdgeDetection(filename)
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
		filename := r.FormValue("image")
		err := app.CreateFlipLR(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to flip image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image flipped successfully, check storage folder for output image")
	})
}

func (app *Config) HandleFlipUD() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateFlipUD(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to flip image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image flipped successfully, check storage folder for output image")
	})
}

func (app *Config) HandleRotate90() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateRotate90(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to rotate image: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image rotated successfully, check storage folder for output image")
	})
}

func (app *Config) HandleRotate270() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("image")
		err := app.CreateRotate270(filename)
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
