//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/dnn"
)

var (
	edgeModelName = "edge_detection_dexined_2024sep"
	edgeModelInit bool
)

// loadConfig loads the configuration for caption size and color from the config store.
// If the configuration is not set, it uses default values.
func loadConfig() {
	if !edgeModelInit {
		logging.Info("Loading model " + edgeModelName)
		net, _, isErr := dnn.NetRead(edgeModelName, "").Result()
		if isErr {
			logging.Error("error loading model " + edgeModelName)
			return
		}
		edgeModelInit = true
		edgeNet = net
	}
}
