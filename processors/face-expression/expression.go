//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-data/face"
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/dnn"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	loadConfig()

	fs := datastore.NewFrameStore(1)
	check := fs.Exists(uint32(image))

	if check.IsErr() || !check.IsOK() {
		logging.Info("no faces for frame")
		return image
	}

	out := image.Clone()
	faces := fs.GetKeys(uint32(image)).Slice()

	for _, f := range faces {
		val, _, isErr := fs.Get(uint32(image), f).Result()
		if isErr {
			logging.Error("Error getting value: " + f)
			return out
		}

		data := face.Data{}
		_, err := data.Read([]byte(val))
		if err != nil {
			logging.Error("Error reading face data: " + err.Error())
			return out
		}

		if emptyRect(data.Rect) {
			logging.Error("Empty face rect")
			continue
		}
		desc := infer(image, data)
		if desc == "" {
			logging.Error("Error inferring face expression")
			continue
		}
		logging.Info("Face expression: " + desc)
		if desc == "unknown" {
			logging.Warn("Unknown face expression")
			continue
		}

		cv.Rectangle(out, data.Rect, green, 1)
		cv.PutText(out, desc, types.Size{X: data.Rect.Min.X, Y: data.Rect.Min.Y - 10}, types.HersheyFontTypeHersheyFontComplex, 0.6, green, 2)
	}

	logging.Debug("Performed face expression analysis on image")

	return out
}

func emptyRect(r types.Rect) bool {
	return r.Min.X == 0 && r.Min.Y == 0 && r.Max.X == 0 && r.Max.Y == 0
}

func stdPoints() mat.Mat {
	pts := mat.MatNewWithSize(5, 1, 13)
	pts.SetFloatAt(0, 0, 38.2946)
	pts.SetFloatAt(0, 1, 51.6963)
	pts.SetFloatAt(1, 0, 73.5318)
	pts.SetFloatAt(1, 1, 51.5014)
	pts.SetFloatAt(2, 0, 56.0252)
	pts.SetFloatAt(2, 1, 71.7366)
	pts.SetFloatAt(3, 0, 41.5493)
	pts.SetFloatAt(3, 1, 92.3655)
	pts.SetFloatAt(4, 0, 70.7299)
	pts.SetFloatAt(4, 1, 92.2041)

	return pts
}

var (
	expressionNet dnn.Net

	expressions = []string{
		"angry", "disgust",
		"fearful", "happy",
		"neutral", "sad", "surprised",
	}

	patchSize         = types.Size{X: 112, Y: 112}
	imageMean float32 = 0.5
	imageStd  float32 = 0.5

	inputName  = "data"
	outputName = "label"

	green = types.RGBA{R: 0, G: 255, B: 0, A: 0}
)

func preprocess(image mat.Mat, points mat.Mat) mat.Mat {
	// image alignment
	transformation, _, isErr := cv.EstimateAffine2d(points, stdPoints()).Result()
	if isErr {
		logging.Error("Error estimating affine transformation")
		return image
	}
	aligned, _, isErr := cv.WarpAffine(image, transformation, patchSize).Result()
	if isErr {
		logging.Error("Error warping image")
		return image
	}

	// image normalization
	normalized, _, isErr := aligned.ConvertToWithParams(mat.MattypeCv32f, 1.0/255.0, 0).Result()
	if isErr {
		logging.Error("Error converting image to float")
		return image
	}
	defer aligned.Close()
	defer normalized.Close()

	normalized.SubtractFloat(imageMean)
	normalized.DivideFloat(imageStd)

	result, _, isErr := dnn.BlobFromImage(normalized, 1.0, types.Size{}, types.Scalar{}, false, false).Result()
	if isErr {
		logging.Error("Error creating blob from image")
		return image
	}

	return result
}

func infer(image mat.Mat, faceData face.Data) string {
	facePoints := convertFacePoints(faceData)
	if facePoints.Empty() {
		logging.Error("Error converting face points")
		return ""
	}
	defer facePoints.Close()

	// create a blob from the image
	blob := preprocess(image, facePoints)
	defer blob.Close()

	// set the input to the network
	expressionNet.SetInput(blob, inputName)

	// run the forward pass
	prob, _, isErr := expressionNet.Forward(outputName).Result()
	if isErr {
		logging.Error("Error running forward pass")
		return ""
	}
	defer prob.Close()

	locs, _, isErr := prob.MinMaxLoc().Result()
	if isErr {
		logging.Error("Error getting min max loc")
		return ""
	}

	return getDescription(float32(locs.MaxLoc.X))
}

func getDescription(index float32) string {
	if index < 0 || int(index) >= len(expressions) {
		return "unknown"
	}

	return expressions[int(index)]
}

func convertFacePoints(data face.Data) mat.Mat {
	// convert facedata to facepoints
	points := mat.MatNewWithSize(5, 1, 13)
	points.SetFloatAt(0, 0, float32(data.RightEye.X))
	points.SetFloatAt(0, 1, float32(data.RightEye.Y))
	points.SetFloatAt(1, 0, float32(data.LeftEye.X))
	points.SetFloatAt(1, 1, float32(data.LeftEye.Y))
	points.SetFloatAt(2, 0, float32(data.NoseTip.X))
	points.SetFloatAt(2, 1, float32(data.NoseTip.Y))
	points.SetFloatAt(3, 0, float32(data.RightMouthCorner.X))
	points.SetFloatAt(3, 1, float32(data.RightMouthCorner.Y))
	points.SetFloatAt(4, 0, float32(data.LeftMouthCorner.X))
	points.SetFloatAt(4, 1, float32(data.LeftMouthCorner.Y))

	return points
}
