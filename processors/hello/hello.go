package main

import (
	"strconv"

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
