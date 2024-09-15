//go:build tinygo

package main

import (
	"github.com/hybridgroup/mechanoid/convert"
	"github.com/wasmvision/wasmcv/components/tinygo/wasm/cv/cv"
	"github.com/wasmvision/wasmcv/components/tinygo/wasm/cv/mat"
	"github.com/wasmvision/wasmcv/components/tinygo/wasm/cv/types"
)

//go:wasmimport hosted println
func println(ptr, size uint32)

//export process
func process(image mat.Mat) mat.Mat {
	imageOut := cv.Blur(image, types.Size{5, 5})
	println(convert.StringToWasmPtr("Performed Blur on image"))

	return imageOut
}

func main() {}
