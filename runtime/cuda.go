//go:build cuda

package runtime

import (
	"log/slog"

	"gocv.io/x/gocv/cuda"
)

func CheckCUDA() bool {
	devices := cuda.GetCudaEnabledDeviceCount()
	if devices == 0 {
		slog.Warn("No CUDA devices found")
		return false
	}
	slog.Info("CUDA devices found", slog.Int("count", devices))
	return true
}
