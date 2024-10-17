package capture

import (
	"github.com/wasmvision/wasmvision/cv"

	"gocv.io/x/gocv"
)

type GStreamer struct {
	device  string
	stream  *gocv.VideoCapture
	retries int
}

func NewGStreamer(device string) *GStreamer {
	return &GStreamer{
		device:  device,
		retries: defaultRetries,
	}
}

func (g *GStreamer) Open() error {
	stream, err := gocv.OpenVideoCaptureWithAPI(g.device, gocv.VideoCaptureGstreamer)
	if err != nil {
		return err
	}

	g.stream = stream
	return nil
}

func (g *GStreamer) Close() error {
	return g.stream.Close()
}

func (g *GStreamer) Read() (*cv.Frame, error) {
	img := gocv.NewMat()
	if ok := g.stream.Read(&img); !ok {
		g.retries--
		if g.retries == 0 {
			return &cv.Frame{}, ErrClosed
		}

		frame := cv.NewEmptyFrame()

		return frame, nil
	}

	g.retries = defaultRetries

	frame := cv.NewFrame(img)

	return frame, nil
}
