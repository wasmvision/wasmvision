//go:build tinygo

package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
)

const defaultCountFrequency = 10 * time.Second

var (
	countFrequency time.Duration
)

// loadConfig loads the configuration for caption size and color from the config store.
// If the configuration is not set, it uses default values.
func loadConfig() {
	if countFrequency == 0 {
		ok, _, isErr := config.GetConfig("face-counter-frequency").Result()
		if isErr {
			countFrequency = defaultCountFrequency
		} else {
			freq, err := strconv.Atoi(strings.Clone(ok))
			if err != nil {
				logging.Error("Error parsing face counter frequency: " + err.Error())
				countFrequency = defaultCountFrequency
			} else {
				countFrequency = time.Duration(freq) * time.Second
			}
		}
		logging.Info("Using face counter frequency seconds: " + defaultCountFrequency.String())
	}
}
