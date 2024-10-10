//go:build tinygo

package main

import (
	"unsafe"

	"github.com/hybridgroup/mechanoid/convert"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/mat"
)

//export process
func process(image mat.Mat) mat.Mat {
	logging.Log("Cols: " +
		convert.IntToString(int(image.Cols())) +
		" Rows: " +
		convert.IntToString(int(image.Rows())) +
		" Type: " +
		convert.IntToString(int(image.Mattype())) +
		" Size: " +
		convert.IntToString(int(image.Size().Len())))

	return image
}

// malloc is needed for wasm-unknown-unknown target for functions that return a List.
//
//export malloc
func malloc(size uint32) uint32 {
	data := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))

	return uint32(ptr)
}

func main() {}
