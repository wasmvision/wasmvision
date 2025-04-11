package capture

import (
	"errors"
	"fmt"

	"github.com/wasmvision/wasmvision/cv"
	"gocv.io/x/gocv"
)

const (
	defaultRetries          = 3
	VideoCaptureFrameHeight = int32(gocv.VideoCaptureFrameHeight)
	VideoCaptureFrameWidth  = int32(gocv.VideoCaptureFrameWidth)
	VideoCaptureFPS         = int32(gocv.VideoCaptureFPS)
)

var ErrClosed = errors.New("capture device closed")

// Capture is the interface that wraps the basic methods for capturing frames.
type Capture interface {
	Open() error
	Close() error
	Read() (*cv.Frame, error)
	Get(property int32) (value float32, err error)
	Set(property int32, value float32) error
}

// OpenDevice opens a capture device based on the specified type and source.
// It returns a Capture interface that can be used to read frames from the device.
// The captureDevice parameter can be "auto", "webcam", "gstreamer", or "ffmpeg".
// The source parameter is the input source for the capture device.
// If the capture device type is not recognized, an error is returned.
// If the capture device fails to open, an error is returned.
func OpenDevice(captureDevice, source string) (Capture, error) {
	var device Capture

	switch captureDevice {
	case "auto", "webcam":
		device = NewWebcam(source)
		if err := device.Open(); err != nil {
			return nil, fmt.Errorf("failed opening video capture: %w", err)
		}
	case "gstreamer":
		device = NewGStreamer(source)
		if err := device.Open(); err != nil {
			return nil, fmt.Errorf("failed opening video capture stream: %w", err)
		}
	case "ffmpeg":
		device = NewFFmpeg(source)
		if err := device.Open(); err != nil {
			return nil, fmt.Errorf("failed opening video capture stream: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown capture type %v", captureDevice)
	}

	return device, nil
}
