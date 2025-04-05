//go:build tinygo

package main

import (
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

var (
	red    = types.RGBA{R: 255, G: 0, B: 0, A: 0}
	green  = types.RGBA{R: 0, G: 255, B: 0, A: 0}
	blue   = types.RGBA{R: 0, G: 0, B: 255, A: 0}
	yellow = types.RGBA{R: 255, G: 255, B: 0, A: 1}
	pink   = types.RGBA{R: 255, G: 105, B: 180, A: 0}
)

// drawRects draws rectangles on the image for the detected objects.
// It uses the class IDs to determine the label for each rectangle.
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
