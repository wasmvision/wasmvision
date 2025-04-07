//go:build tinygo

package main

import (
	"strconv"

	"github.com/wasmvision/wasmvision-data/face"
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
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

	facedata [60]byte
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

	loadConfig()

	sz := image.Size().Slice()
	detector.SetInputSize(types.Size{X: int32(sz[1]), Y: int32(sz[0])})

	faces, _, isErr := detector.Detect(image).Result()
	if isErr {
		logging.Error("error detecting faces")
		return image
	}

	defer faces.Close()

	out := image
	if drawFaceBoxes {
		out = image.Clone()
	}
	handleFaces(out, faces)

	logging.Info("Performed face detection on image " + strconv.Itoa(int(uint32(out))))

	return out
}

func handleFaces(out mat.Mat, faces mat.Mat) {
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

		if drawFaceBoxes {
			cv.Rectangle(out, faceRect, green, 1)
			cv.Circle(out, rightEye, 1, blue, 1)
			cv.Circle(out, leftEye, 1, red, 1)
			cv.Circle(out, noseTip, 1, green, 1)
			cv.Circle(out, rightMouthCorner, 1, pink, 1)
			cv.Circle(out, leftMouthCorner, 1, yellow, 1)
		}
	}
}

func storeFaceData(image mat.Mat, faceid int, faceRect types.Rect, rightEye, leftEye, noseTip, rightMouthCorner, leftMouthCorner types.Size) {
	face := face.Data{
		ID:               uint32(faceid),
		Rect:             faceRect,
		RightEye:         rightEye,
		LeftEye:          leftEye,
		NoseTip:          noseTip,
		RightMouthCorner: rightMouthCorner,
		LeftMouthCorner:  leftMouthCorner,
	}

	if _, err := face.Write(facedata[:]); err != nil {
		logging.Error("error writing face data: " + err.Error())
		return
	}
	fid := "face-" + strconv.Itoa(faceid+1)
	_, _, isErr := fs.Set(uint32(image), fid, string(facedata[:])).Result()
	if isErr {
		logging.Error("error setting value")
	}
}
