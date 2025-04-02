//go:build tinygo

package main

import (
	"encoding/binary"

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
		logging.Info("no faces for frame")
		return out
	}

	faces := fs.GetKeys(uint32(image)).Slice()

	for _, face := range faces {
		val, _, isErr := fs.Get(uint32(image), face).Result()
		if isErr {
			logging.Error("Error getting value: " + face)
			return out
		}

		rect := faceRect([]byte(val))
		if emptyRect(rect) {
			logging.Error("empty rect for face")
			continue
		}

		area := out.Region(rect)
		blurred, _, isErr := cv.Blur(area, types.Size{X: 50, Y: 50}).Result()
		if isErr {
			logging.Error("Error applying blur")
			return out
		}

		blurred.CopyTo(area)

		blurred.Close()
		area.Close()
	}

	logging.Info("Performed face blur on image")

	return out
}

// faceRect returns a types.Rect from a byte slice of face data.
func faceRect(data []byte) (faceRect types.Rect) {
	if len(data) < 16 {
		logging.Error("faceRect: invalid data length")
		return
	}

	return types.Rect{
		Min: types.Size{
			X: int32(binary.LittleEndian.Uint32(data[0:4])),
			Y: int32(binary.LittleEndian.Uint32(data[4:8])),
		},
		Max: types.Size{
			X: int32(binary.LittleEndian.Uint32(data[8:12])),
			Y: int32(binary.LittleEndian.Uint32(data[12:16])),
		},
	}
}

func emptyRect(r types.Rect) bool {
	return r.Min.X == 0 && r.Min.Y == 0 && r.Max.X == 0 && r.Max.Y == 0
}
