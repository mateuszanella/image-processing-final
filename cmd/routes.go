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
	mux.Handle("POST /api/image-combination", app.HandleCombinationUploadImage())
	mux.Handle("GET /api/image", app.HandleGetImage())
	mux.Handle("GET /api/image-combination", app.HandleGetCombinationImage())
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
	mux.Handle("POST /api/min-sdf", app.HandleMinimumFilter())
	mux.Handle("POST /api/max-sdf", app.HandleMaximumFilter())
	mux.Handle("POST /api/order-sdf", app.HandleOrderFilter())
	mux.Handle("POST /api/conservative-smoothing-sdf", app.HandleConservativeSmoothingFilter())

	mux.Handle("POST /api/dilation", app.HandleDilation())
	mux.Handle("POST /api/erosion", app.HandleErosion())
	mux.Handle("POST /api/opening", app.HandleOpening())
	mux.Handle("POST /api/closing", app.HandleClosing())
	mux.Handle("POST /api/contour", app.HandleContour())

	mux.Handle("POST /api/prewitt", app.HandlePrewittEdgeDetection())
	mux.Handle("POST /api/sobel", app.HandleSobelEdgeDetection())
	mux.Handle("POST /api/laplacian", app.HandleLaplacianEdgeDetection())

	mux.Handle("POST /api/flip-lr", app.HandleFlipLR())
	mux.Handle("POST /api/flip-ud", app.HandleFlipUD())
	mux.Handle("POST /api/rotate-90", app.HandleRotate90())
	mux.Handle("POST /api/rotate-270", app.HandleRotate270())

	// mux.Handle(("POST /api/combination/add"), app.HandleAddImages())

	// templ routes
	mux.Handle("/", templ.Handler(layout.Base(view.Index())))
	mux.Handle("/combination", templ.Handler(layout.Base(view.CombinationPage())))
	mux.Handle("/blank", templ.Handler(partials.Blank()))
	mux.Handle("GET /image", app.HandleDisplayComponent())
	mux.Handle("GET /component/dropzone", app.HandleDropzoneComponent())
	mux.Handle("GET /component/combination-dropzone", app.HandleCombinationDropzoneComponent())
	mux.Handle("GET /component/filters", app.HandleFiltersComponent())
	mux.Handle("GET /component/combination-filters", app.HandleCombinationFiltersComponent())
	mux.Handle("GET /component/adjustments", app.HandleAdjustmentsComponent())

	return mux
}
