//go:build tinygo

package main

import (
	"strings"

	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"github.com/wasmvision/wasmvision-sdk-go/vlm"
)

var (
	modelName     string
	projectorName string
	prompt        string
	modelInit     bool
)

const (
	defaultModel     = "Qwen2.5-VL-3B-Instruct-Q8_0"
	defaultProjector = "mmproj-Qwen2.5-VL-3B-Instruct-Q8_0"
	defaultPrompt    = "Describe what is in this picture in highly complimentary terms using 6 words or less."
)

func loadConfig() {
	if modelName == "" {
		ok, _, isErr := config.GetConfig("vlm-model").Result()
		if isErr {
			modelName = defaultModel
		} else {
			modelName = strings.Clone(ok)
		}

		logging.Info("Using VLM model " + modelName)
	}

	if projectorName == "" {
		ok, _, isErr := config.GetConfig("vlm-projector").Result()
		if isErr {
			projectorName = defaultProjector
		} else {
			projectorName = strings.Clone(ok)
		}

		logging.Info("Using VLM projector " + projectorName)
	}

	if prompt == "" {
		ok, _, isErr := config.GetConfig("vlm-prompt").Result()
		if isErr {
			prompt = defaultPrompt
		} else {
			prompt = strings.Clone(ok)
		}

		logging.Info("Using prompt " + prompt)
	}

	if !modelInit {
		logging.Info("Loading VLM " + modelName)
		m, _, isErr := vlm.ModelInitFromFile(modelName, projectorName).Result()
		if isErr {
			logging.Error("Error loading VLM model " + modelName)
			return
		}
		model = m
		modelInit = true
	}
}
