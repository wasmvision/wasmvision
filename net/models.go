package net

type ModelFile struct {
	Alias    string
	Filename string
	URL      string
}

var KnownModels = map[string]ModelFile{
	"mosaic-9": {
		Alias:    "mosaic-9",
		Filename: "mosaic-9.onnx",
		URL:      "https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-9.onnx",
	},
	"mosaic-8": {
		Alias:    "mosaic-8",
		Filename: "mosaic-8.onnx",
		URL:      "https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-8.onnx",
	},
}
