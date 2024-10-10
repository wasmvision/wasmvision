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

	red    = types.RGBA{R: 255, G: 0, B: 0, A: 0}
	green  = types.RGBA{R: 0, G: 255, B: 0, A: 0}
	blue   = types.RGBA{R: 0, G: 0, B: 255, A: 0}
	yellow = types.RGBA{R: 255, G: 255, B: 0, A: 1}
	pink   = types.RGBA{R: 255, G: 105, B: 180, A: 0}
)

func init() {
	detector = objdetect.NewFaceDetectorYN("face_detection_yunet_2023mar", "", types.Size{X: 200, Y: 200})
}

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
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

		rightEye := types.Size{
			X: int32(faces.GetFloatAt(r, 4)),
			Y: int32(faces.GetFloatAt(r, 5)),
		}

		leftEye := types.Size{
			X: int32(faces.GetFloatAt(r, 6)),
			Y: int32(faces.GetFloatAt(r, 7)),
		}

		noseTip := types.Size{
			X: int32(faces.GetFloatAt(r, 8)),
			Y: int32(faces.GetFloatAt(r, 9)),
		}

		rightMouthCorner := types.Size{
			X: int32(faces.GetFloatAt(r, 10)),
			Y: int32(faces.GetFloatAt(r, 11)),
		}

		leftMouthCorner := types.Size{
			X: int32(faces.GetFloatAt(r, 12)),
			Y: int32(faces.GetFloatAt(r, 13)),
		}

		cv.Rectangle(out, faceRect, green, 1)
		cv.Circle(out, rightEye, 1, blue, 1)
		cv.Circle(out, leftEye, 1, red, 1)
		cv.Circle(out, noseTip, 1, green, 1)
		cv.Circle(out, rightMouthCorner, 1, pink, 1)
		cv.Circle(out, leftMouthCorner, 1, yellow, 1)
	}

	logging.Log("Performed face detection on image")

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
