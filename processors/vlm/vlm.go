//go:build tinygo

package main

import (
	"time"

	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	hosttime "github.com/wasmvision/wasmvision-sdk-go/time"
	"github.com/wasmvision/wasmvision-sdk-go/vlm"
	"wasmcv.org/wasm/cv/mat"
)

var (
	lastUpdate time.Time
)

var model vlm.Model

func init() {
	lastUpdate = time.UnixMicro(int64(hosttime.Now(0)))
}

//export process
func process(image mat.Mat) mat.Mat {
	loadConfig()

	now := time.UnixMicro(int64(hosttime.Now(0)))
	if now.Sub(lastUpdate) > 5*time.Second {
		logging.Info("Asking for image description...")

		data, _, isErr := model.Prompt(prompt, uint32(image)).Result()
		switch {
		case isErr:
			logging.Error("VLM error")
		case len(data) > 0:
			ps := datastore.NewProcessorStore(1)
			ps.Set("captions", "caption", data)
			logging.Info(data)
		default:
			logging.Info("No result from VLM")
		}

		lastUpdate = now
	}

	return image
}
