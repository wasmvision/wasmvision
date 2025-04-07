//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/config"
)

var (
	checkedDrawFaceBoxes bool
	drawFaceBoxes        bool
)

func loadConfig() {
	if !checkedDrawFaceBoxes {
		drawFaceBoxes = true
		checkedDrawFaceBoxes = true

		ok, _, isErr := config.GetConfig("detect-draw-faces").Result()
		if !isErr && ok == "false" {
			drawFaceBoxes = false
		}
	}
}
