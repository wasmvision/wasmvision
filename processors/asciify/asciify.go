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
	resized := cv.Resize(image, types.Size{X: 80, Y: 60}, 0, 0, types.InterpolationTypeInterpolationLinear)

	imageToAscii(resized)

	// output to terminal
	for y := 0; y < 60; y++ {
		println(convert.StringToWasmPtr(string(ascii[y][:])))
	}

	return resized
}

func main() {}
