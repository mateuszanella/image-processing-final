package main

import (
	"fmt"
	"image-processing/view/partials"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

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

		fileName := uuid.New().String() + ".jpg"

		out, err := os.Create("./storage/" + fileName)
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

		file.Seek(0, 0)

		out2, err := os.Create("./storage/" + "output.jpg")
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
		filePath := "./storage/output.jpg"

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
		thresholdStr := r.FormValue("threshold")
		threshold := 128
		var err error
		if thresholdStr != "" {
			threshold, err = strconv.Atoi(thresholdStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse threshold: %v", err), http.StatusInternalServerError)
				return
			}
		}

		err = app.CreateBinary(filename, uint8(threshold))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create grayscale: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created sucessfully, check storage folder for output image")
	})
}

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
