//go:build tinygo

package main

import (
	"github.com/hybridgroup/mechanoid/convert"
	"github.com/wasmvision/wasmcv/components/tinygo/wasm/cv/mat"
)

//go:wasmimport hosted println
func println(ptr, size uint32)

//export process
func process(image mat.Mat) mat.Mat {
	println(convert.StringToWasmPtr("Cols: " +
		convert.IntToString(int(image.Cols())) +
		" Rows: " +
		convert.IntToString(int(image.Rows())) +
		" Type: " +
		convert.IntToString(int(image.Type()))))

	return image
}

func main() {}
