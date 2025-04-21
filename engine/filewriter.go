package engine

import (
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/cv"
)

type FileWriter interface {
	Close()
	Write(*cv.Frame) error
	Start(capture.Capture) error
}
