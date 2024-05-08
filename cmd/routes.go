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

	mux.Handle("POST /api/add", app.HandleAddValue())
	mux.Handle("POST /api/subtract", app.HandleSubtractValue())
	mux.Handle("POST /api/multiply", app.HandleMultiplyValue())
	mux.Handle("POST /api/divide", app.HandleDivideValue())

	mux.Handle("POST /api/not", app.HandleNotOpertion())

	mux.Handle("POST /api/negative", app.HandleCreateNegativeFilter())
	mux.Handle("POST /api/histogram-equalization", app.HandleCreateHisogramEqualization())

	mux.Handle("POST /api/mean-sdf", app.HandleMeanFilter())
	mux.Handle("POST /api/median-sdf", app.HandleMedianFilter())
	mux.Handle("POST /api/gaussian-sdf", app.HandleGaussianFilter())

	mux.Handle("POST /api/dilation", app.HandleDilation())
	mux.Handle("POST /api/erosion", app.HandleErosion())
	mux.Handle("POST /api/opening", app.HandleOpening())
	mux.Handle("POST /api/closing", app.HandleClosing())
	mux.Handle("POST /api/contour", app.HandleContour())

	mux.Handle("POST /api/prewitt", app.HandlePrewittEdgeDetection())

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
