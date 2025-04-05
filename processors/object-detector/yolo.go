//go:build tinygo

package main

import (
	"strconv"
	"strings"

	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"go.bytecodealliance.org/cm"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/dnn"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

// Array of YOLOv8 class labels
var classes = []string{
	"person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck", "boat",
	"traffic light", "fire hydrant", "stop sign", "parking meter", "bench", "bird", "cat", "dog", "horse",
	"sheep", "cow", "elephant", "bear", "zebra", "giraffe", "backpack", "umbrella", "handbag", "tie",
	"suitcase", "frisbee", "skis", "snowboard", "sports ball", "kite", "baseball bat", "baseball glove",
	"skateboard", "surfboard", "tennis racket", "bottle", "wine glass", "cup", "fork", "knife", "spoon",
	"bowl", "banana", "apple", "sandwich", "orange", "broccoli", "carrot", "hot dog", "pizza", "donut",
	"cake", "chair", "couch", "potted plant", "bed", "dining table", "toilet", "tv", "laptop", "mouse",
	"remote", "keyboard", "cell phone", "microwave", "oven", "toaster", "sink", "refrigerator", "book",
	"clock", "vase", "scissors", "teddy bear", "hair drier", "toothbrush",
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
			logging.Error("Error getting layer")
			continue
		}
		layerName, _, isErr := layer.GetName().Result()
		if isErr {
			logging.Error("Error getting layer name")
			continue
		}
		if layerName != "_input" {
			outputLayers = append(outputLayers, strings.Clone(layerName))
		}
	}

	return outputLayers
}

var (
	yoloNet dnn.Net

	red    = types.RGBA{R: 255, G: 0, B: 0, A: 0}
	green  = types.RGBA{R: 0, G: 255, B: 0, A: 0}
	blue   = types.RGBA{R: 0, G: 0, B: 255, A: 0}
	yellow = types.RGBA{R: 255, G: 255, B: 0, A: 1}
	pink   = types.RGBA{R: 255, G: 105, B: 180, A: 0}

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

func detect(src mat.Mat) ([]dnn.Rect, []int, []uint32) {
	blob, _, isErr := dnn.BlobFromImageWithParams(src, params).Result()
	if isErr {
		logging.Error("Error creating blob from image")
		return []dnn.Rect{}, []int{}, []uint32{}
	}
	defer blob.Close()

	if blob.Empty() {
		logging.Error("Blob is empty")
		return []dnn.Rect{}, []int{}, []uint32{}
	}
	// feed the blob into the detector
	yoloNet.SetInput(blob, "")

	// run a forward pass thru the network
	probs, _, isErr := yoloNet.ForwardLayers(cm.ToList(outputNames)).Result()
	if isErr {
		logging.Error("Error running forward pass")
		return []dnn.Rect{}, []int{}, []uint32{}
	}

	ps := make([]mat.Mat, len(probs.Slice()))
	copy(ps, probs.Slice())

	defer func() {
		logging.Debug("Closing probs" + strconv.Itoa(len(ps)))
		for _, prob := range ps {
			prob.Close()
		}
	}()

	boxes, confidences, classIds := performDetection(ps)
	if len(boxes) == 0 {
		return []dnn.Rect{}, []int{}, []uint32{}
	}

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

	bs := make([]dnn.Rect, len(iboxes.Slice()))
	copy(bs, iboxes.Slice())

	return bs, classIds, indices.Slice()
}

func performDetection(outs []mat.Mat) ([]types.Rect, []float32, []int) {
	var classIds []int
	var confidences []float32
	var boxes []types.Rect

	// needed for yolov8
	var isErr bool
	original := outs[0]
	defer original.Close()

	outs[0], _, isErr = cv.TransposeND(original, cm.ToList[[]int32]([]int32{0, 2, 1})).Result()
	if isErr {
		logging.Error("Error transposing output")
		return boxes, confidences, classIds
	}

	for _, out := range outs {
		out, _, isErr = out.Reshape(1, out.Size().Slice()[1]).Result()
		if isErr {
			logging.Error("Error reshaping output")
			return boxes, confidences, classIds
		}

		cols := out.Cols()
		for i := uint32(0); i < out.Rows(); i++ {
			scoresCol, _, isErr := out.RowRange(i, i+1).Result()
			if isErr {
				logging.Error("Error getting scores column")
				return boxes, confidences, classIds
			}
			if scoresCol.Empty() {
				logging.Error("Scores column is empty")
				return boxes, confidences, classIds
			}

			scores, _, _isErr := scoresCol.ColRange(4, cols).Result()
			if _isErr {
				logging.Error("Error getting scores")
				return boxes, confidences, classIds
			}
			if scores.Empty() {
				logging.Error("Scores are empty")
				return boxes, confidences, classIds
			}

			result, _, isErr := scores.MinMaxLoc().Result()
			if isErr {
				logging.Error("Error getting min max loc")
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

func drawRects(img mat.Mat, boxes []types.Rect, classIds []int, indices []uint32) {
	for _, idx := range indices {
		if idx == 0 {
			continue
		}
		rect := types.Rect{Min: types.Size{X: boxes[idx].Min.X, Y: boxes[idx].Min.Y}, Max: types.Size{X: boxes[idx].Max.X, Y: boxes[idx].Max.Y}}
		cv.Rectangle(img, rect, green, 2)
		cv.PutText(img, classes[classIds[idx]], types.Size{X: boxes[idx].Min.X, Y: boxes[idx].Min.Y - 10}, types.HersheyFontTypeHersheyFontComplex, 0.6, green, 2)
	}
}
