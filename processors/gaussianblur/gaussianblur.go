//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//export process
func process(image mat.Mat) mat.Mat {
	imageOut := cv.GaussianBlur(image, types.Size{X: 25, Y: 25}, 4.5, 4.5, types.BorderTypeBorderReflect101)
	logging.Info("Performed GaussianBlur on image")

	return imageOut
}
