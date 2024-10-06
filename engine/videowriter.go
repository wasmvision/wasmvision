package engine

import (
	"log"

	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/runtime"
	"gocv.io/x/gocv"
)

const (
	defaultCodec    = "MJPG"
	defaultFBS      = 25
	framebufferSize = 3
)

// VideoWriter represents a file writer used for saving video frames
type VideoWriter struct {
	writer   *gocv.VideoWriter
	Filename string
	codec    string
	fps      float64
	refs     *runtime.MapRefs
	frames   chan *cv.Frame
}

func NewVideoWriter(refs *runtime.MapRefs, dest string) VideoWriter {
	return VideoWriter{
		Filename: dest,
		codec:    defaultCodec,
		fps:      defaultFBS,
		refs:     refs,
		frames:   make(chan *cv.Frame, framebufferSize),
	}
}

func (vw *VideoWriter) Close() {
	vw.writer.Close()
}

func (vw *VideoWriter) Write(img *cv.Frame) error {
	vw.frames <- img
	return nil
}

func (vw *VideoWriter) Start(source capture.Capture) error {
	sample, err := source.Read()
	if err != nil {
		return err
	}

	defer sample.Close()

	videoWriter, err := gocv.VideoWriterFile(vw.Filename, vw.codec, vw.fps, sample.Image.Cols(), sample.Image.Rows(), true)
	if err != nil {
		return err
	}

	go vw.writeFrames()

	vw.writer = videoWriter

	return nil
}

func (vw *VideoWriter) writeFrames() {
	for frame := range vw.frames {
		if err := vw.writer.Write(frame.Image); err != nil {
			log.Printf("error writing frame: %v\n", err)
		}

		frame.Close()
		vw.refs.Drop(frame.ID.Unwrap())
	}
}
