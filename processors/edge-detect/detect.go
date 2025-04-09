//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/dnn"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

var (
	edgeNet dnn.Net
)

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	loadConfig()

	blob, _, isErr := dnn.BlobFromImage(image, 1.0, types.Size{X: 512, Y: 512}, types.Scalar{Val1: float32(103.5), Val2: float32(116.2), Val3: float32(123.6)}, false, false).Result()
	if isErr {
		logging.Error("error creating blob")
		return image
	}
	defer blob.Close()

	edgeNet.SetInput(blob, "")
	result := analyze(image)

	return result
}

func analyze(image mat.Mat) mat.Mat {
	result, _, isErr := edgeNet.Forward("").Result()
	if isErr {
		logging.Error("error running model")
		return image
	}

	data := []mat.Mat{result}
	i := len(data) - 1
	if data[i].Empty() {
		logging.Warn("output was empty")
		return image
	}
	var processed mat.Mat
	processed, _, isErr = data[i].Reshape(0, 512).Result()
	if isErr {
		logging.Error("error reshaping output")
		return image
	}
	defer processed.Close()

	inp := sigmoid(processed)
	defer inp.Close()

	normalized, _, isErr := cv.Normalize(inp, 0, 255, uint32(32)).Result()
	if isErr {
		logging.Error("error normalizing output")
		return image
	}
	defer normalized.Close()

	cvt, _, isErr := normalized.ConvertTo(mat.MattypeCv8u).Result()
	if isErr {
		logging.Error("error converting output")
		return image
	}

	resized, _, isErr := cv.Resize(cvt, types.Size{X: int32(image.Cols()), Y: int32(image.Rows())}, 0, 0, types.InterpolationTypeInterpolationLinear).Result()
	if isErr {
		logging.Error("error resizing output")
		return image
	}

	return resized
}

// Function to apply sigmoid activation
func sigmoid(input mat.Mat) mat.Mat {
	// e^-input
	input.MultiplyFloat(-1.0)
	inp, _, isErr := cv.Exp(input).Result()
	if isErr {
		logging.Error("error applying exp")
		return input
	}
	defer inp.Close()

	// (1 + e^-input)
	inp.AddFloat(1.0)

	// 1 / (1 + e^-input)
	ones, _, isErr := mat.MatOnes(input.Rows(), input.Cols(), mat.MattypeCv32f).Result()
	if isErr {
		logging.Error("error creating ones matrix")
		return input
	}
	defer ones.Close()

	out, _, isErr := cv.Divide(ones, inp).Result()
	if isErr {
		logging.Error("error applying divide")
		return input
	}

	return out
}
