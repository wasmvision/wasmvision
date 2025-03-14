//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

var captionText []byte

func init() {
	captionText = make([]byte, 0, 80)
}

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

	l := len(result.OK().Slice())
	switch {
	case l == 0:
		logging.Info("empty caption")
		return out
	case l > 80:
		l = 80
	}

	captionText = append(captionText[:0], result.OK().Slice()[:l]...)

	var msg string
	msg = string(captionText)
	logging.Info("caption: " + msg)
	cv.PutText(out, msg, types.Point{X: 10, Y: 20}, types.HersheyFontTypeHersheyFontSimplex, 0.6, types.RGBA{R: 255, G: 255, B: 255, A: 0}, 2)

	return out
}
