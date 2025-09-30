//go:build !llama

package runtime

import (
	"log/slog"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
)

func handleLlama(modules wypes.Modules, enable bool, cctx *cv.Context) error {
	slog.Warn("you cannot enable llama.cpp in this build of wasmVision")
	return nil
}

func llamaFree() {
}
