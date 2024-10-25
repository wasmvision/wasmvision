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

var candyNet dnn.Net

func init() {
	candyNet = dnn.NetRead("candy-9", "")
}

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() || candyNet.Empty() {
		logging.Warn("image was empty")
		return image
	}

	// convert image Mat to 320x240 blob that the style transfer can analyze
	blob := dnn.BlobFromImage(image, 1.0,
		types.Size{X: 320, Y: 240},
		types.Scalar{Val1: redAdjust, Val2: greenAdjust, Val3: blueAdjust, Val4: 0},
		false, false)
	defer blob.Close()

	// feed the blob into the detector
	candyNet.SetInput(blob, "")

	probMat := candyNet.Forward("")
	sz := probMat.Size().Slice()
	dims := sz[2] * sz[3]

	candied := mat.MatNewWithSize(240, 320, 16)
	defer candied.Close()

	// take blob and obtain displayable Mat image from it
	for i := uint32(0); i < dims; i++ {
		r := probMat.GetFloatAt(0, i)
		r += redAdjust

		g := probMat.GetFloatAt(0, i+dims)
		g += greenAdjust

		b := probMat.GetFloatAt(0, i+dims*2)
		b += blueAdjust

		candied.SetUcharAt(0, i*3, uint8(r))
		candied.SetUcharAt(0, i*3+1, uint8(g))
		candied.SetUcharAt(0, i*3+2, uint8(b))
	}

	// resize back to original size
	out := cv.Resize(candied, types.Size{X: int32(image.Cols()), Y: int32(image.Rows())}, 0, 0, types.InterpolationTypeInterpolationLinear)

	logging.Info("Performed neural style transfer on image")

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
