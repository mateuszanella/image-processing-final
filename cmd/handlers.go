package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

		fmt.Fprintln(w, "Image uploaded successfully")
	})
}

func (app *Config) HandleGetImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir("./storage")
		if err != nil {
			http.Error(w, "Failed to read storage", http.StatusInternalServerError)
			return
		}

		if len(files) == 0 {
			http.Error(w, "No images found", http.StatusNotFound)
			return
		}

		file := files[len(files)-1]

		http.ServeFile(w, r, "./storage/"+file.Name())
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
		filename := r.FormValue("image")
		err := app.CreateGrayscale(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create grayscale: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image created sucessfully, check storage folder for output image")
	})
}
