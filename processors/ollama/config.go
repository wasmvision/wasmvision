//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
)

var (
	url      string
	model    string
	prompt   string
	template string
)

const (
	defaultURL    = "http://localhost:11434"
	defaultModel  = "llava"
	defaultPrompt = "Describe what is in this picture in highly complimentary terms using 6 words or less."
)

func loadConfig() {
	if url == "" {
		ok, _, isErr := config.GetConfig("ollama-url").Result()
		if isErr {
			url = defaultURL
		} else {
			url = ok
		}

		logging.Info("Using Ollama server at " + url)
	}

	if model == "" {
		ok, _, isErr := config.GetConfig("ollama-model").Result()
		if isErr {
			model = defaultModel
		} else {
			model = ok
		}

		logging.Info("Using Ollama model " + model)
	}

	if prompt == "" {
		ok, _, isErr := config.GetConfig("ollama-prompt").Result()
		if isErr {
			prompt = defaultPrompt
		} else {
			prompt = ok
		}

		logging.Info("Using prompt " + prompt)
	}

	if template == "" {
		template = `{
			"model": "` + model + `",
			"prompt":"` + prompt + `",
			"stream": false,
			"images": ["%IMAGE%"]
		  }`
	}
}
