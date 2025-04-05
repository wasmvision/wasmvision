//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/dnn"
)

var (
	yoloModelName string
	yoloModelInit bool
)

// loadConfig loads the configuration for caption size and color from the config store.
// If the configuration is not set, it uses default values.
func loadConfig() {
	if yoloModelName == "" {
		ok, _, isErr := config.GetConfig("yolo-model").Result()
		if isErr {
			yoloModelName = "yolov8n"
		} else {
			yoloModelName = ok
		}

		logging.Info("Using YOLO model " + yoloModelName)
	}

	if !yoloModelInit {
		logging.Info("Loading model " + yoloModelName)
		net, _, isErr := dnn.NetRead(yoloModelName, "").Result()
		if isErr {
			logging.Error("Error loading model " + yoloModelName)
			return
		}
		yoloModelInit = true
		yoloNet = net

		outputNames = getOutputNames(yoloNet)
	}
}
