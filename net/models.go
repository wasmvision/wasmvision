package net

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
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
}

func DownloadModel(name string, modelsDir string) error {
	model, ok := KnownModels[name]
	if !ok {
		return errors.New("unknown model")
	}

	opts := []getter.ClientOption{}
	client := &getter.Client{
		Ctx:     context.Background(),
		Src:     model.URL,
		Dst:     filepath.Join(modelsDir, filepath.Base(model.Filename)),
		Mode:    getter.ClientModeFile,
		Options: opts,
	}

	if err := client.Get(); err != nil {
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
