//go:build tinygo

package main

import (
	"strconv"
	"strings"

	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//export process
func process(image mat.Mat) mat.Mat {
	loadConfig()

	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	out := image.Clone()

	ps := datastore.NewProcessorStore(1)
	check := ps.Exists("captions")

	if check.IsErr() || !check.IsOK() {
		logging.Info("no captions")
		return out
	}

	ok, err, isErr := ps.Get("captions", "caption").Result()
	if isErr {
		logging.Info("no caption found: " + strconv.Itoa(int(err)))
		return out
	}

	captions := wrapCaption(ok, int(captionWordsPerLine))
	if len(captions) == 0 {
		logging.Info("empty caption")
		return out
	}

	y := int32(captionLineHeight)
	for _, caption := range captions {
		cv.PutText(out, caption, types.Point{X: 10, Y: y}, types.HersheyFontTypeHersheyFontSimplex, captionSize, captionColor, 2)
		logging.Info(caption)
		y += int32(captionLineHeight)
	}

	return out
}

// wrapCaption splits a string into multiple lines with a specified word limit per line.
// It returns a slice of strings, each representing a line of the caption.
// If the input string is empty or contains only whitespace, it returns a slice with the original string.
func wrapCaption(s string, limit int) []string {
	if strings.TrimSpace(s) == "" {
		return []string{s}
	}

	words := strings.Fields(strings.Clone(s))
	var result []string
	start := 0

	for start < len(words) {
		end := start + limit
		if end > len(words) {
			end = len(words)
		}
		result = append(result, strings.Join(words[start:end], " "))
		start = end
	}

	return result
}
