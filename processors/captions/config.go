//go:build tinygo

package main

import (
	"strconv"
	"strings"

	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/types"
)

const (
	captionWordsPerLineDefault = 10
	captionSizeDefault         = 0.6
	captionLineHeightDefault   = 20
)

var (
	captionWordsPerLine int64
	captionLineHeight   int64
	captionSize         float64
	captionColorName    string
	captionColor        types.RGBA

	captionColorDefault = captionColorBlack
	captionColorBlack   = types.RGBA{R: 0, G: 0, B: 0, A: 0}
	captionColorWhite   = types.RGBA{R: 255, G: 255, B: 255, A: 0}
	captionColorRed     = types.RGBA{R: 255, G: 0, B: 0, A: 0}
	captionColorGreen   = types.RGBA{R: 0, G: 255, B: 0, A: 0}
	captionColorBlue    = types.RGBA{R: 0, G: 0, B: 255, A: 0}
	captionColorYellow  = types.RGBA{R: 255, G: 255, B: 0, A: 0}
)

// loadConfig loads the configuration for caption size and color from the config store.
// If the configuration is not set, it uses default values.
func loadConfig() {
	if captionWordsPerLine == 0 {
		ok, _, isErr := config.GetConfig("caption-words-per-line").Result()
		if isErr {
			captionWordsPerLine = captionWordsPerLineDefault
		} else {
			wpl, err := strconv.Atoi(strings.Clone(ok))
			if err != nil {
				logging.Error("Error parsing words per line: " + err.Error())
				captionWordsPerLine = captionWordsPerLineDefault
			} else {
				captionWordsPerLine = int64(wpl)
			}
		}
		logging.Info("Using caption words per line " + strconv.Itoa(int(captionWordsPerLine)))
	}

	if captionLineHeight == 0 {
		ok, _, isErr := config.GetConfig("caption-line-height").Result()
		if isErr {
			captionLineHeight = captionLineHeightDefault
		} else {
			ht, err := strconv.Atoi(strings.Clone(ok))
			if err != nil {
				logging.Error("Error parsing line height: " + err.Error())
				captionLineHeight = captionLineHeightDefault
			} else {
				captionLineHeight = int64(ht)
			}
		}

		logging.Info("Using caption line height " + strconv.Itoa(int(captionLineHeight)))
	}

	if captionSize == 0 {
		ok, _, isErr := config.GetConfig("caption-size").Result()
		if isErr {
			captionSize = captionSizeDefault
		} else {
			captionSize, _ = strconv.ParseFloat(strings.Clone(ok), 64)
		}

		logging.Info("Using caption size " + strconv.FormatFloat(captionSize, 'f', -1, 64))
	}

	if captionColorName == "" {
		ok, _, isErr := config.GetConfig("caption-color").Result()
		if isErr {
			captionColorName = "black"
			captionColor = captionColorBlack
		} else {
			captionColorName = ok
			switch captionColorName {
			case "black":
				captionColor = captionColorBlack
			case "white":
				captionColor = captionColorWhite
			case "red":
				captionColor = captionColorRed
			case "green":
				captionColor = captionColorGreen
			case "blue":
				captionColor = captionColorBlue
			case "yellow":
				captionColor = captionColorYellow
			default:
				captionColorName = "black"
				// default to black if the color is not recognized
				captionColor = captionColorBlack
			}
		}

		logging.Info("Using color " + captionColorName)
	}
}
