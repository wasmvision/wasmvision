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
	// start with a blank screen
	println(convert.StringToWasmPtr("\033[2J\033[3J\033[H"))

	resized := cv.Resize(image, types.Size{X: 80, Y: 60}, 0, 0, types.InterpolationTypeInterpolationNearest)
	defer resized.Close()

	imageToAscii(resized)

	// output to terminal
	for y := 0; y < 60; y++ {
		println(convert.StringToWasmPtr(string(ascii[y][:])))
	}

	asciified := mat.MatNewMatWithSize(image.Rows(), image.Cols(), 16)
	for y := 0; y < 60; y++ {
		cv.PutText(asciified, string(ascii[y][:]), types.Size{X: 10, Y: 8 * int32(y)}, types.HersheyFontTypeHersheyFontComplexSmall, 0.8, types.Rgba{R: 255, G: 255, B: 255, A: 255}, 1)
	}

	return asciified
}

func main() {}
