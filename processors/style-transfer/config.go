//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/dnn"
)

var (
	styleModelName string
	styleModelInit bool
)

// loadConfig loads the configuration for caption size and color from the config store.
// If the configuration is not set, it uses default values.
func loadConfig() {
	if styleModelName == "" {
		ok, _, isErr := config.GetConfig("style-model").Result()
		if isErr {
			styleModelName = "mosaic-9"
		} else {
			styleModelName = ok
		}

		logging.Info("Using style model " + styleModelName)
	}

	if !styleModelInit {
		logging.Info("Loading style model " + styleModelName)
		net, _, isErr := dnn.NetRead(styleModelName, "").Result()
		if isErr {
			logging.Error("Error loading style model " + styleModelName)
			return
		}
		styleNet = net
		styleModelInit = true
	}
}
