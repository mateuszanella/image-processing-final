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
	mux.Handle("GET /api/test", app.HandleTestImageManipulation())

	mux.Handle("POST /api/grayscale", app.HandleCreateGrayscale())
	mux.Handle("POST /api/binary", app.HandleCreateBinary())

	// templ routes
	c := layout.Base(view.Index())
	mux.Handle("/", templ.Handler(c))
	mux.Handle("/blank", templ.Handler(partials.Blank()))
	mux.Handle("/foo", templ.Handler(partials.Foo()))
	mux.Handle("GET /image", app.HandleDisplayComponent())
	mux.Handle("GET /component/dropzone", app.HandleDropzoneComponent())
	mux.Handle("GET /component/filters", app.HandleFiltersComponent())

	return mux
}
