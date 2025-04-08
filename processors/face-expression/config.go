//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/dnn"
)

var (
	modelInit bool
	modelName = "face_expression_recognition_mobilefacenet_2022july"
)

func loadConfig() {
	if !modelInit {
		logging.Info("Loading model " + modelName)
		net, _, isErr := dnn.NetRead(modelName, "").Result()
		if isErr {
			logging.Error("Error loading model " + modelName)
			return
		}
		modelInit = true
		expressionNet = net
	}
}
