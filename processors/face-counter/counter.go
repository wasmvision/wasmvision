//go:build tinygo

package main

import (
	"time"

	"github.com/orsinium-labs/jsony"
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	hosttime "github.com/wasmvision/wasmvision-sdk-go/time"
	"wasmcv.org/wasm/cv/mat"
)

var (
	lastUpdate time.Time

	facesCount int
	frameCount int
)

func init() {
	lastUpdate = time.UnixMicro(int64(hosttime.Now(0)))
}

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	loadConfig()

	fs := datastore.NewFrameStore(1)
	check := fs.Exists(uint32(image))

	if check.IsErr() || !check.IsOK() {
		logging.Info("no faces for frame")
		return image
	}

	faces := fs.GetKeys(uint32(image)).Slice()
	facesCount += len(faces)
	frameCount++

	now := time.UnixMicro(int64(hosttime.Now(0)))
	if now.Sub(lastUpdate) > countFrequency {
		avg := float32(facesCount) / float32(frameCount)
		facesCount = 0
		frameCount = 0
		lastUpdate = now

		ps := datastore.NewProcessorStore(1)
		tm := now.Format(time.RFC3339Nano)
		obj := jsony.Object{
			jsony.Field{"timestamp", jsony.String(tm)},
			jsony.Field{"average-faces-seen", jsony.Float32(avg)},
		}
		s := jsony.EncodeString(obj)
		ps.Set("face-counter", tm, s)

		logging.Info(s)
	}

	return image
}
