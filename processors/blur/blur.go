//go:build tinygo

package main

import (
	"github.com/hybridgroup/mechanoid/convert"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//go:wasmimport hosted println
func println(ptr, size uint32)

//export process
func process(image mat.Mat) mat.Mat {
	imageOut := cv.Blur(image, types.Size{X: 25, Y: 25})
	println(convert.StringToWasmPtr("Performed Blur on image"))

	return imageOut
}

func main() {}
