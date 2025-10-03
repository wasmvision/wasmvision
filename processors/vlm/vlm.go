//go:build tinygo

package main

import (
	"strings"
	"time"

	"github.com/orsinium-labs/jsony"
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
			result := strings.Clone(data)
			storeCaption(result)
			logging.Info(result)
		default:
			logging.Info("No result from VLM")
		}

		lastUpdate = now
	}

	return image
}

func storeCaption(caption string) {
	ps := datastore.NewProcessorStore(1)
	now := time.UnixMicro(int64(hosttime.Now(0)))
	tm := now.Format(time.RFC3339Nano)
	obj := jsony.Object{
		jsony.Field{"timestamp", jsony.String(tm)},
		jsony.Field{"caption", jsony.String(caption)},
	}
	data := jsony.EncodeString(obj)
	ps.Set("captions", "caption", data)
}
