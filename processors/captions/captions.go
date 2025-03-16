//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	out := image.Clone()

	ps := datastore.NewProcessorStore(1)
	check := ps.Exists("captions")

	if check.IsErr() || !check.IsOK() {
		logging.Info("no captions")
		return out
	}

	result := ps.Get("captions", "caption")
	if result.IsErr() {
		logging.Info("no caption")
		return out
	}

	formatCaption(*result.OK())

	if caption == "" {
		logging.Info("empty caption")
		return out
	}

	cv.PutText(out, caption, types.Point{X: 10, Y: 30}, types.HersheyFontTypeHersheyFontSimplex, 1.0, types.RGBA{R: 0, G: 0, B: 0, A: 0}, 2)
	logging.Info("caption: " + caption)

	return out
}

var (
	caption string
)

func formatCaption(msg string) {
	switch {
	case len(msg) == 0:
		caption = ""
	case len(msg) < 32:
		caption = string(msg)
	default:
		caption = string(msg[:32]) + "..."
	}
}
