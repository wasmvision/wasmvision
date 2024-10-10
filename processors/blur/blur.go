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
	imageOut := cv.Blur(image, types.Size{X: 25, Y: 25})
	logging.Log("Performed Blur on image")

	return imageOut
}

func main() {}
