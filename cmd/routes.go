package main

import (
	"net/http"

	"image-processing/view"
	"image-processing/view/layout"
	"image-processing/view/partials"

	"github.com/a-h/templ"
)

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", app.HandleStaticFiles())

	mux.Handle("POST /api/image", app.HandleUploadImage())
	mux.Handle("GET /api/image", app.HandleGetImage())
	mux.Handle("GET /api/image/{id}", app.HandleGetImageByID())

	// templ routes
	c := layout.Base(view.Index())
	mux.Handle("/", templ.Handler(c))
	mux.Handle("/foo", templ.Handler(partials.Foo()))
	mux.Handle("GET /image/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "No image ID provided", http.StatusBadRequest)
			return
		}

		c := partials.ImageDisplay(id)

		templ.Handler(c).ServeHTTP(w, r)
	}))

	return mux
}
