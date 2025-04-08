package models

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter/v2"
)

type ModelFile struct {
	Alias    string
	Filename string
	URL      string
}

var KnownModels = map[string]ModelFile{
	"candy-9": {
		Alias:    "candy-9",
		Filename: "candy-9.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/candy-9.onnx",
	},
	"candy-8": {
		Alias:    "candy-8",
		Filename: "candy-8.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/candy-8.onnx",
	},
	"mosaic-9": {
		Alias:    "mosaic-9",
		Filename: "mosaic-9.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-9.onnx",
	},
	"mosaic-8": {
		Alias:    "mosaic-8",
		Filename: "mosaic-8.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-8.onnx",
	},
	"pointilism-9": {
		Alias:    "pointilism-9",
		Filename: "pointilism-9.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/pointilism-9.onnx",
	},
	"pointilism-8": {
		Alias:    "pointilism-8",
		Filename: "pointilism-8.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/pointilism-8.onnx",
	},
	"rain-princess-9": {
		Alias:    "rain-princess-9",
		Filename: "rain-princess-9.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/rain-princess-9.onnx",
	},
	"rain-princess-8": {
		Alias:    "rain-princess-8",
		Filename: "rain-princess-8.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/rain-princess-8.onnx",
	},
	"udnie-9": {
		Alias:    "udnie-9",
		Filename: "udnie-9.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/udnie-9.onnx",
	},
	"udnie-8": {
		Alias:    "udnie-8",
		Filename: "udnie-8.onnx",
		URL:      "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/udnie-8.onnx",
	},
	"face_detection_yunet_2023mar": {
		Alias:    "face_detection_yunet_2023mar",
		Filename: "face_detection_yunet_2023mar.onnx",
		URL:      "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/face_detection_yunet/face_detection_yunet_2023mar.onnx",
	},
	"yolov8n": {
		Alias:    "yolov8n",
		Filename: "yolov8n.onnx",
		URL:      "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8n.onnx",
	},
	"yolov8s": {
		Alias:    "yolov8s",
		Filename: "yolov8s.onnx",
		URL:      "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8s.onnx",
	},
	"yolov8m": {
		Alias:    "yolov8m",
		Filename: "yolov8m.onnx",
		URL:      "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8m.onnx",
	},
	"yolov8l": {
		Alias:    "yolov8l",
		Filename: "yolov8l.onnx",
		URL:      "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8l.onnx",
	},
	"yolov8x": {
		Alias:    "yolov8x",
		Filename: "yolov8x.onnx",
		URL:      "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8x.onnx",
	},
	"yolox": {
		Alias:    "yolox",
		Filename: "object_detection_yolox_2022nov.onnx",
		URL:      "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/object_detection_yolox/object_detection_yolox_2022nov.onnx",
	},
	"face_expression_recognition_mobilefacenet_2022july": {
		Alias:    "face_expression_recognition_mobilefacenet_2022july",
		Filename: "face_expression_recognition_mobilefacenet_2022july.onnx",
		URL:      "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july.onnx",
	},
	"facial_expression_recognition_mobilefacenet_2022july_int8": {
		Alias:    "facial_expression_recognition_mobilefacenet_2022july_int8",
		Filename: "facial_expression_recognition_mobilefacenet_2022july_int8.onnx",
		URL:      "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july_int8.onnx",
	},
	"facial_expression_recognition_mobilefacenet_2022july_int8bq": {
		Alias:    "facial_expression_recognition_mobilefacenet_2022july_int8bq",
		Filename: "facial_expression_recognition_mobilefacenet_2022july_int8bq.onnx",
		URL:      "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july_int8bq.onnx",
	},
}

func Download(name string, modelsDir string) error {
	model, ok := KnownModels[name]
	if !ok {
		return errors.New("unknown model")
	}

	req := &getter.Request{
		Src:     model.URL,
		Dst:     filepath.Join(modelsDir, filepath.Base(model.Filename)),
		GetMode: getter.ModeFile,
	}

	client := &getter.Client{}
	if _, err := client.Get(context.Background(), req); err != nil {
		return err
	}

	return nil
}

// ModelFile gets the model file path name for the Net.
func ModelFileName(model string, modelsDir string) string {
	if km, ok := KnownModels[model]; ok {
		return filepath.Join(modelsDir, km.Filename)
	}

	return filepath.Join(modelsDir, model)
}

func ModelExists(model string) bool {
	if _, err := os.Stat(model); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func ModelWellKnown(model string) bool {
	if _, ok := KnownModels[model]; ok {
		return true
	}

	return false
}
