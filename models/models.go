package models

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter/v2"
)

type ModelFile struct {
	Alias       string
	Filename    string
	URL         string
	Description string
}

var KnownModels = map[string]ModelFile{
	"candy-9": {
		Alias:       "candy-9",
		Filename:    "candy-9.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/candy-9.onnx",
		Description: "Candy model for fast neural style transfer (opset 8)",
	},
	"candy-8": {
		Alias:       "candy-8",
		Filename:    "candy-8.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/candy-8.onnx",
		Description: "Candy model for fast neural style transfer (opset 8)",
	},
	"mosaic-9": {
		Alias:       "mosaic-9",
		Filename:    "mosaic-9.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-9.onnx",
		Description: "Mosaic model for fast neural style transfer (opset 9)",
	},
	"mosaic-8": {
		Alias:       "mosaic-8",
		Filename:    "mosaic-8.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-8.onnx",
		Description: "Mosaic model for fast neural style transfer (opset 8)",
	},
	"pointilism-9": {
		Alias:       "pointilism-9",
		Filename:    "pointilism-9.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/pointilism-9.onnx",
		Description: "Pointilism model for fast neural style transfer (opset 9)",
	},
	"pointilism-8": {
		Alias:       "pointilism-8",
		Filename:    "pointilism-8.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/pointilism-8.onnx",
		Description: "Pointilism model for fast neural style transfer (opset 8)",
	},
	"rain-princess-9": {
		Alias:       "rain-princess-9",
		Filename:    "rain-princess-9.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/rain-princess-9.onnx",
		Description: "Rain Princess model for fast neural style transfer (opset 9)",
	},
	"rain-princess-8": {
		Alias:       "rain-princess-8",
		Filename:    "rain-princess-8.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/rain-princess-8.onnx",
		Description: "Rain Princess model for fast neural style transfer (opset 8)",
	},
	"udnie-9": {
		Alias:       "udnie-9",
		Filename:    "udnie-9.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/udnie-9.onnx",
		Description: "Udnie model for fast neural style transfer (opset 9)",
	},
	"udnie-8": {
		Alias:       "udnie-8",
		Filename:    "udnie-8.onnx",
		URL:         "https://github.com/onnx/models/raw/refs/heads/main/validated/vision/style_transfer/fast_neural_style/model/udnie-8.onnx",
		Description: "Udnie model for fast neural style transfer (opset 8)",
	},
	"yunet_2023mar": {
		Alias:       "yunet_2023mar",
		Filename:    "face_detection_yunet_2023mar.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/face_detection_yunet/face_detection_yunet_2023mar.onnx",
		Description: "Yunet face detection model",
	},
	"yolov8n": {
		Alias:       "yolov8n",
		Filename:    "yolov8n.onnx",
		URL:         "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8n.onnx",
		Description: "YOLOv8n model for object detection",
	},
	"yolov8s": {
		Alias:       "yolov8s",
		Filename:    "yolov8s.onnx",
		URL:         "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8s.onnx",
		Description: "YOLOv8s model for object detection",
	},
	"yolov8m": {
		Alias:       "yolov8m",
		Filename:    "yolov8m.onnx",
		URL:         "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8m.onnx",
		Description: "YOLOv8m model for object detection",
	},
	"yolov8l": {
		Alias:       "yolov8l",
		Filename:    "yolov8l.onnx",
		URL:         "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8l.onnx",
		Description: "YOLOv8l model for object detection",
	},
	"yolov8x": {
		Alias:       "yolov8x",
		Filename:    "yolov8x.onnx",
		URL:         "https://huggingface.co/cabelo/yolov8/resolve/main/yolov8x.onnx",
		Description: "YOLOv8x model for object detection",
	},
	"yolox": {
		Alias:       "yolox",
		Filename:    "object_detection_yolox_2022nov.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/object_detection_yolox/object_detection_yolox_2022nov.onnx",
		Description: "YOLOX model for object detection",
	},
	"mobilefacenet_2022july": {
		Alias:       "mobilefacenet_2022july",
		Filename:    "face_expression_recognition_mobilefacenet_2022july.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july.onnx",
		Description: "MobileFaceNet model for Facial Expression Recognition (FER)",
	},
	"mobilefacenet_2022july_int8": {
		Alias:       "mobilefacenet_2022july_int8",
		Filename:    "facial_expression_recognition_mobilefacenet_2022july_int8.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july_int8.onnx",
		Description: "MobileFaceNet model for Facial Expression Recognition (quantized)",
	},
	"mobilefacenet_2022july_int8bq": {
		Alias:       "mobilefacenet_2022july_int8bq",
		Filename:    "facial_expression_recognition_mobilefacenet_2022july_int8bq.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/facial_expression_recognition/facial_expression_recognition_mobilefacenet_2022july_int8bq.onnx",
		Description: "MobileFaceNet model for Facial Expression Recognition (quantized with bias)",
	},
	"dexined_2024sep": {
		Alias:       "dexined_2024sep",
		Filename:    "edge_detection_dexined_2024sep.onnx",
		URL:         "https://github.com/opencv/opencv_zoo/raw/refs/heads/main/models/edge_detection_dexined/edge_detection_dexined_2024sep.onnx",
		Description: "Dexined model for edge detection",
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
