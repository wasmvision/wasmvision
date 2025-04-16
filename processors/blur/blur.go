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
	out, _, isErr := cv.Blur(image, types.Size{X: 25, Y: 25}).Result()
	if isErr {
		logging.Error("Error applying blur")
		return image
	}
	logging.Debug("Performed Blur on image")

	if out.Empty() {
		logging.Warn("Blurred image was empty")
		return image
	}
	return out
}
