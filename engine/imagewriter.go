package engine

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/runtime"
	"gocv.io/x/gocv"
)

// ImageWriter represents a file writer used for saving image frames
type ImageWriter struct {
	Filename string
	refs     *runtime.MapRefs
	frames   chan *cv.Frame
}

func NewImageWriter(refs *runtime.MapRefs, dest string) *ImageWriter {
	return &ImageWriter{
		Filename: dest,
		refs:     refs,
		frames:   make(chan *cv.Frame, framebufferSize),
	}
}

func (iw *ImageWriter) Close() {
	time.Sleep(500 * time.Millisecond)

	close(iw.frames)
}

func (iw *ImageWriter) Write(img *cv.Frame) error {
	iw.frames <- img
	return nil
}

func (iw *ImageWriter) Start(source capture.Capture) error {
	go iw.writeFrames()

	return nil
}

func (iw *ImageWriter) writeFrames() {
	for frame := range iw.frames {
		filename := strings.Replace(iw.Filename, "{id}", strconv.Itoa(int(frame.ID.Unwrap())), -1)
		if !gocv.IMWrite(filename, frame.Image) {
			slog.Error("error writing frame", "id", strconv.Itoa(int(frame.ID.Unwrap())))
		}

		slog.Info("wrote frame", "id", strconv.Itoa(int(frame.ID.Unwrap())), "to", filename)

		frame.Close()
		iw.refs.Drop(frame.ID.Unwrap())
	}
}
