package engine

import (
	"log"

	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/frame"
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
	cache    *frame.Cache
	frames   chan frame.Frame
}

func NewVideoWriter(cache *frame.Cache, dest string) VideoWriter {
	return VideoWriter{
		Filename: dest,
		codec:    defaultCodec,
		fps:      defaultFBS,
		cache:    cache,
		frames:   make(chan frame.Frame, framebufferSize),
	}
}

func (vw *VideoWriter) Close() {
	vw.writer.Close()
}

func (vw *VideoWriter) Write(img frame.Frame) error {
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
		vw.cache.Delete(frame.ID)
	}
}
