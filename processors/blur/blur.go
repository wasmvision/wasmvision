//go:build tinygo

package main

import (
	"strconv"

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
	logging.Info("Performed Blur on image")

	logging.Info("Blurred image: " + strconv.Itoa(int(image)) + " " + strconv.Itoa(int(out)))
	if out.Empty() {
		logging.Warn("Blurred image was empty")
		return image
	}
	return out
}
