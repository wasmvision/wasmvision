//go:build tinygo

package main

import (
	"github.com/wasmvision/wasmvision-data/face"
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	"wasmcv.org/wasm/cv/cv"
	"wasmcv.org/wasm/cv/mat"
	"wasmcv.org/wasm/cv/types"
)

//export process
func process(image mat.Mat) mat.Mat {
	if image.Empty() {
		logging.Warn("image was empty")
		return image
	}

	out := image.Clone()

	fs := datastore.NewFrameStore(1)
	check := fs.Exists(uint32(image))

	if check.IsErr() || !check.IsOK() {
		logging.Warn("no faces for frame")
		return out
	}

	faces := fs.GetKeys(uint32(image)).Slice()

	for _, f := range faces {
		val, _, isErr := fs.Get(uint32(image), f).Result()
		if isErr {
			logging.Error("Error getting value: " + f)
			return out
		}

		data := face.Data{}
		_, err := data.Read([]byte(val))
		if err != nil {
			logging.Error("Error reading face data: " + err.Error())
			return out
		}

		area := out.Region(data.Rect)
		blurred, _, isErr := cv.Blur(area, types.Size{X: 50, Y: 50}).Result()
		if isErr {
			logging.Error("Error applying blur")
			return out
		}

		blurred.CopyTo(area)

		blurred.Close()
		area.Close()
	}

	logging.Debug("Performed face blur on image")

	return out
}

func emptyRect(r types.Rect) bool {
	return r.Min.X == 0 && r.Min.Y == 0 && r.Max.X == 0 && r.Max.Y == 0
}
