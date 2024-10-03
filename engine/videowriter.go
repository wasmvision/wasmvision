package engine

import (
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/frame"
	"gocv.io/x/gocv"
)

const (
	defaultCodec = "MJPG"
	defaultFBS   = 25
)

// VideoWriter represents a file writer used for saving video frames
type VideoWriter struct {
	writer   *gocv.VideoWriter
	Filename string
	codec    string
	fps      float64
}

func NewVideoWriter(dest string) VideoWriter {
	return VideoWriter{
		Filename: dest,
		codec:    defaultCodec,
		fps:      defaultFBS,
	}
}

func (vw *VideoWriter) Close() {
	vw.writer.Close()
}

func (vw *VideoWriter) Write(img frame.Frame) error {
	return vw.writer.Write(img.Image)
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

	vw.writer = videoWriter

	return nil
}
