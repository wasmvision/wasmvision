package main

import (
	"strconv"
	"unsafe"

	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/mat"
)

//export process
func process(image mat.Mat) mat.Mat {
	logging.Info("Cols: " +
		strconv.Itoa(int(image.Cols())) +
		" Rows: " +
		strconv.Itoa(int(image.Rows())) +
		" Type: " +
		strconv.Itoa(int(image.Mattype())) +
		" Size: " +
		strconv.Itoa(int(image.Size().Len())))

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
