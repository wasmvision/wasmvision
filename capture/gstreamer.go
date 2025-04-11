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

func (g *GStreamer) Get(property int32) (value float32, err error) {
	return float32(g.stream.Get(gocv.VideoCaptureProperties(property))), nil
}

func (g *GStreamer) Set(property int32, value float32) error {
	g.stream.Set(gocv.VideoCaptureProperties(property), float64(value))
	return nil
}
