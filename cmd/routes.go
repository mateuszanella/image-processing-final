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

	mux.Handle("POST /api/upload", app.HandleUploadImage())
	mux.Handle("GET /api/image", app.HandleGetImage())
	mux.Handle("GET /api/test", app.HandleTestImageManipulation())

	// templ routes
	c := layout.Base(view.Index())
	mux.Handle("/", templ.Handler(c))
	mux.Handle("/foo", templ.Handler(partials.Foo()))
	mux.Handle("GET /image", templ.Handler(partials.ImageDisplay()))

	return mux
}
