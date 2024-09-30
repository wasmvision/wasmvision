//go:build tinygo

package main

import (
	"unsafe"

	"github.com/hybridgroup/mechanoid/convert"
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

//go:wasmimport hosted println
func println(ptr, size uint32)

var mosaicNet dnn.Net

func init() {
	mosaicNet = dnn.NetReadNet("mosaic-9.onnx", "")
}

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() || mosaicNet.Empty() {
		return mat.MatNewMat()
	}

	// convert image Mat to 320x240 blob that the style transfer can analyze
	blob := dnn.BlobFromImage(image, 1.0,
		types.Size{X: 320, Y: 240},
		types.Scalar{Val1: redAdjust, Val2: greenAdjust, Val3: blueAdjust, Val4: 0},
		false, false)
	defer blob.Close()

	// feed the blob into the detector
	mosaicNet.SetInput(blob, "")

	probMat := mosaicNet.Forward("")
	sz := probMat.Size().Slice()
	dims := sz[2] * sz[3]

	mosaiced := mat.MatNewMatWithSize(240, 320, 16)
	defer mosaiced.Close()

	// take blob and obtain displayable Mat image from it
	for i := uint32(0); i < dims; i++ {
		r := probMat.GetFloatAt(0, i)
		r += redAdjust

		g := probMat.GetFloatAt(0, i+dims)
		g += greenAdjust

		b := probMat.GetFloatAt(0, i+dims*2)
		b += blueAdjust

		mosaiced.SetUcharAt(0, i*3, uint8(r))
		mosaiced.SetUcharAt(0, i*3+1, uint8(g))
		mosaiced.SetUcharAt(0, i*3+2, uint8(b))
	}

	// resize back to original size
	out := cv.Resize(mosaiced, types.Size{X: int32(image.Cols()), Y: int32(image.Rows())}, 0, 0, types.InterpolationTypeInterpolationLinear)

	println(convert.StringToWasmPtr("Performed neural style transfer on image"))

	return out
}

func main() {}

// malloc is needed for wasm-unknown-unknown target for functions that return a List.
//
//export malloc
func malloc(size uint32) uint32 {
	data := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))

	return uint32(ptr)
}
