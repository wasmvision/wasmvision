//go:build tinygo

package main

import (
	"unsafe"

	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/dnn"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

const (
	redAdjust   = 103.939
	greenAdjust = 116.779
	blueAdjust  = 123.68
)

var styleNet dnn.Net

//export process
func process(image mat.Mat) mat.Mat {
	loadConfig()

	switch {
	case image.Empty():
		logging.Warn("image was empty")
		return image

	case styleNet.Empty():
		logging.Warn("style DNN was empty")
		return image
	}

	// convert image Mat to 320x240 blob that the style transfer can analyze
	blob := dnn.BlobFromImage(image, 1.0,
		types.Size{X: 320, Y: 240},
		types.Scalar{Val1: redAdjust, Val2: greenAdjust, Val3: blueAdjust, Val4: 0},
		false, false)
	defer blob.Close()

	// feed the blob into the detector
	styleNet.SetInput(blob, "")

	// perform the style transfer
	// and get the result
	// the result is a blob with 3 channels
	// and 240x320 size
	results := styleNet.Forward("")
	defer results.Close()

	sz := results.Size().Slice()
	dims := sz[2] * sz[3]

	styled := mat.MatNewWithSize(240, 320, 16)
	defer styled.Close()

	// take blob and obtain displayable Mat image from it
	// by drawing the 3 channels on to the styled output
	for i := uint32(0); i < dims; i++ {
		r := results.GetFloatAt(0, i)
		r += redAdjust

		g := results.GetFloatAt(0, i+dims)
		g += greenAdjust

		b := results.GetFloatAt(0, i+dims*2)
		b += blueAdjust

		styled.SetUcharAt(0, i*3, uint8(r))
		styled.SetUcharAt(0, i*3+1, uint8(g))
		styled.SetUcharAt(0, i*3+2, uint8(b))
	}

	// resize the styled output back to original size so it can be displayed
	out := cv.Resize(styled, types.Size{X: int32(image.Cols()), Y: int32(image.Rows())}, 0, 0, types.InterpolationTypeInterpolationLinear)

	logging.Info("Performed neural style transfer on image")

	return out
}

// malloc is needed for wasm-unknown-unknown target for functions that return a List.
//
//export malloc
func malloc(size uint32) uint32 {
	data := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))

	return uint32(ptr)
}
