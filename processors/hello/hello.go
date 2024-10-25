package main

import (
	"strconv"
	"unsafe"

	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/mat"
)

var configValue string

//export process
func process(image mat.Mat) mat.Mat {
	if configValue == "" {
		conf := config.GetConfig("default")
		if conf.IsErr() {
			configValue = conf.Err().String()
			logging.Error("Config error: " + configValue)
		} else {
			configValue = *conf.OK()
			logging.Info("Config: " + configValue)
		}
	}

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

func main() {}
