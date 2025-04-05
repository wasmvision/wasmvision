//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/mat"
)

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	loadConfig()

	if len(outputNames) == 0 {
		logging.Warn("No output names provided")
		return image
	}

	boxes, classIds, indices := detection(image)
	if len(boxes) == 0 {
		logging.Info("No classes detected")
		return image
	}

	out := image.Clone()
	drawRects(out, boxes, classIds, indices)

	return out
}
