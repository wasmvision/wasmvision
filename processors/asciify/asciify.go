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
	// start with a blank screen
	logging.Println("\033[2J\033[3J\033[H")

	resized := cv.Resize(image, types.Size{X: 80, Y: 60}, 0, 0, types.InterpolationTypeInterpolationNearest)
	defer resized.Close()

	imageToAscii(resized)

	// output to terminal
	for y := 0; y < 60; y++ {
		logging.Println(string(ascii[y][:]))
	}

	asciified := mat.MatNewWithSize(image.Rows(), image.Cols(), 16)
	for y := 0; y < 60; y++ {
		cv.PutText(asciified, string(ascii[y][:]),
			types.Size{X: 10, Y: 8 * int32(y)}, types.HersheyFontTypeHersheyFontComplexSmall, 0.8,
			types.RGBA{R: 255, G: 255, B: 255, A: 255}, 1)
	}

	return asciified
}

func main() {}
