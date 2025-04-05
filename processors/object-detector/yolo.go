//go:build tinygo

package main

import (
	"strings"

	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"go.bytecodealliance.org/cm"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/dnn"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

var (
	yoloNet dnn.Net

	outputNames []string

	ratio    = float32(0.003921568627)
	mean     = types.Scalar{Val1: 0, Val2: 0, Val3: 0, Val4: 0}
	swapRGB  = false
	padValue = types.Scalar{Val1: 144.0, Val2: 0, Val3: 0, Val4: 0}

	scoreThreshold float32 = 0.5
	nmsThreshold   float32 = 0.4

	params = types.BlobParams{ScaleFactor: ratio,
		Size:        types.Size{X: 640, Y: 640},
		Mean:        mean,
		SwapRB:      swapRGB,
		Ddepth:      uint8(mat.MattypeCv32f),
		DataLayout:  types.DataLayoutTypeDataLayoutNchw,
		PaddingMode: types.PaddingModeTypePaddingModeLetterbox,
		Border:      padValue}
)

// process the image and return the detected objects
func detection(src mat.Mat) ([]dnn.Rect, []int, []uint32) {
	// prepare the image
	blob, _, isErr := dnn.BlobFromImageWithParams(src, params).Result()
	if isErr {
		logging.Error("error creating blob from image")
		return []dnn.Rect{}, []int{}, []uint32{}
	}
	defer blob.Close()

	if blob.Empty() {
		logging.Error("blob is empty")
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	// feed the blob into the detector, and run a forward pass
	// through the network.
	yoloNet.SetInput(blob, "")
	probs, _, isErr := yoloNet.ForwardLayers(cm.ToList(outputNames)).Result()
	if isErr {
		logging.Error("error running forward pass")
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	ps := make([]mat.Mat, len(probs.Slice()))
	copy(ps, probs.Slice())

	defer func() {
		for _, prob := range ps {
			prob.Close()
		}
	}()

	// analyze the output probabilities
	// and get the bounding boxes, confidences, and class IDs
	// for the detected objects
	boxes, confidences, classIds := analyze(ps)
	if len(boxes) == 0 {
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	// calculate tracking info for detected objects
	sz := types.Size{X: int32(src.Cols()), Y: int32(src.Rows())}
	iboxes, _, isErr := dnn.BlobRectsToImageRects(params, cm.ToList(boxes), sz).Result()
	if isErr {
		logging.Error("Error converting boxes to image rects")
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	indices, _, isErr := dnn.NMSBoxes(iboxes, cm.ToList(confidences), scoreThreshold, nmsThreshold).Result()
	if isErr {
		logging.Error("Error running NMS")
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	return iboxes.Slice(), classIds, indices.Slice()
}

// analyze the output of the network
func analyze(outs []mat.Mat) ([]types.Rect, []float32, []int) {
	var boxes []types.Rect
	var confidences []float32
	var classIds []int

	// transpose needed for yolov8
	var isErr bool
	original := outs[0]
	defer original.Close()

	outs[0], _, isErr = cv.TransposeND(original, cm.ToList[[]int32]([]int32{0, 2, 1})).Result()
	if isErr {
		logging.Error("error transposing output")
		return boxes, confidences, classIds
	}

	for _, out := range outs {
		out, _, isErr = out.Reshape(1, out.Size().Slice()[1]).Result()
		if isErr {
			logging.Error("error reshaping output")
			return boxes, confidences, classIds
		}

		cols := out.Cols()
		for i := uint32(0); i < out.Rows(); i++ {
			scoresCol, _, isErr := out.RowRange(i, i+1).Result()
			if isErr {
				logging.Error("error getting scores column")
				return boxes, confidences, classIds
			}
			if scoresCol.Empty() {
				logging.Error("scores column is empty")
				return boxes, confidences, classIds
			}

			scores, _, _isErr := scoresCol.ColRange(4, cols).Result()
			if _isErr {
				logging.Error("error getting scores")
				return boxes, confidences, classIds
			}
			if scores.Empty() {
				logging.Error("scores are empty")
				return boxes, confidences, classIds
			}

			result, _, isErr := scores.MinMaxLoc().Result()
			if isErr {
				logging.Error("error getting min max loc")
				return boxes, confidences, classIds
			}
			confidence, classIDPoint := result.MaxVal, result.MaxLoc

			if confidence > 0.5 {
				centerX := out.GetFloatAt(i, cols)
				centerY := out.GetFloatAt(i, cols+1)
				width := out.GetFloatAt(i, cols+2)
				height := out.GetFloatAt(i, cols+3)

				left := centerX - width/2
				top := centerY - height/2
				right := centerX + width/2
				bottom := centerY + height/2
				classIds = append(classIds, int(classIDPoint.X))
				confidences = append(confidences, float32(confidence))

				rect := types.Rect{Min: types.Size{X: int32(left), Y: int32(top)}, Max: types.Size{X: int32(right), Y: int32(bottom)}}
				boxes = append(boxes, rect)
			}
		}
	}

	return boxes, confidences, classIds
}

func getOutputNames(net dnn.Net) []string {
	var outputLayers []string
	layers, _, isErr := net.GetUnconnectedOutLayers().Result()
	if isErr {
		logging.Error("Error getting unconnected out layers")
		return outputLayers
	}

	for i := 0; i < len(layers.Slice()); i++ {
		layer, _, isErr := net.GetLayer(uint32(layers.Slice()[i])).Result()
		if isErr {
			logging.Error("error getting layer")
			continue
		}
		layerName, _, isErr := layer.GetName().Result()
		if isErr {
			logging.Error("error getting layer name")
			continue
		}
		if layerName != "_input" {
			outputLayers = append(outputLayers, strings.Clone(layerName))
		}
	}

	return outputLayers
}
