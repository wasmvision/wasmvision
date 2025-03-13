//go:build tinygo

package main

import (
	"encoding/binary"
	"strconv"
	"unsafe"

	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"go.bytecodealliance.org/cm"
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

	facedata [8 * 7]byte
	fs       datastore.FrameStore
)

func init() {
	detector = objdetect.NewFaceDetectorYN("face_detection_yunet_2023mar", "", types.Size{X: 200, Y: 200})

	fs = datastore.NewFrameStore(1)
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

		storeFaceData(out, int(r), faceRect, rightEye, leftEye, noseTip, rightMouthCorner, leftMouthCorner)

		cv.Rectangle(out, faceRect, green, 1)
		cv.Circle(out, rightEye, 1, blue, 1)
		cv.Circle(out, leftEye, 1, red, 1)
		cv.Circle(out, noseTip, 1, green, 1)
		cv.Circle(out, rightMouthCorner, 1, pink, 1)
		cv.Circle(out, leftMouthCorner, 1, yellow, 1)
	}

	logging.Info("Performed face detection on image " + strconv.Itoa(int(uint32(out))))

	return out
}

func storeFaceData(image mat.Mat, faceid int, faceRect types.Rect, rightEye, leftEye, noseTip, rightMouthCorner, leftMouthCorner types.Size) {
	binary.LittleEndian.PutUint32(facedata[0:4], uint32(faceRect.Min.X))
	binary.LittleEndian.PutUint32(facedata[4:8], uint32(faceRect.Min.Y))
	binary.LittleEndian.PutUint32(facedata[8:12], uint32(faceRect.Max.X))
	binary.LittleEndian.PutUint32(facedata[12:16], uint32(faceRect.Max.Y))

	binary.LittleEndian.PutUint32(facedata[16:20], uint32(rightEye.X))
	binary.LittleEndian.PutUint32(facedata[20:24], uint32(rightEye.Y))

	binary.LittleEndian.PutUint32(facedata[24:28], uint32(leftEye.X))
	binary.LittleEndian.PutUint32(facedata[28:32], uint32(leftEye.Y))

	binary.LittleEndian.PutUint32(facedata[32:36], uint32(noseTip.X))
	binary.LittleEndian.PutUint32(facedata[36:40], uint32(noseTip.Y))

	binary.LittleEndian.PutUint32(facedata[40:44], uint32(rightMouthCorner.X))
	binary.LittleEndian.PutUint32(facedata[44:48], uint32(rightMouthCorner.Y))

	binary.LittleEndian.PutUint32(facedata[48:52], uint32(leftMouthCorner.X))
	binary.LittleEndian.PutUint32(facedata[52:56], uint32(leftMouthCorner.Y))

	fid := "face-" + strconv.Itoa(faceid+1)
	val := fs.Set(uint32(image), fid, cm.ToList(facedata[:]))
	if val.IsErr() {
		logging.Error("Error setting value: " + val.Err().String())
	}
}

// malloc is needed for wasm-unknown-unknown target for functions that return a List.
//
//export malloc
func malloc(size uint32) uint32 {
	data := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))

	return uint32(ptr)
}
