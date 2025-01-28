//go:build tinygo

package main

import (
	"unsafe"

	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/objdetect"
	"wasmcv.org/wasm/cv/types"
)

var (
	detector objdetect.FaceDetectorYN
)

func init() {
	detector = objdetect.NewFaceDetectorYN("face_detection_yunet_2023mar", "", types.Size{X: 200, Y: 200})
}

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	sz := image.Size().Slice()
	detector.SetInputSize(types.Size{X: int32(sz[1]), Y: int32(sz[0])})

	faces := detector.Detect(image)
	defer faces.Close()

	out := image.Clone()

	for r := uint32(0); r < faces.Rows(); r++ {
		x0 := int32(faces.GetFloatAt(r, 0))
		y0 := int32(faces.GetFloatAt(r, 1))
		x1 := x0 + int32(faces.GetFloatAt(r, 2))
		y1 := y0 + int32(faces.GetFloatAt(r, 3))

		faceRect := types.Rect{Min: types.Size{X: x0, Y: y0}, Max: types.Size{X: x1, Y: y1}}

		area := out.Region(faceRect)
		blurred := cv.Blur(area, types.Size{X: 50, Y: 50})
		blurred.CopyTo(area)

		blurred.Close()
		area.Close()
	}

	logging.Info("Performed face blur on image")

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
