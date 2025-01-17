package capture

import (
	"github.com/wasmvision/wasmvision/cv"

	"gocv.io/x/gocv"
)

type FFmpeg struct {
	device  string
	stream  *gocv.VideoCapture
	retries int
}

func NewFFmpeg(device string) *FFmpeg {
	return &FFmpeg{
		device:  device,
		retries: defaultRetries,
	}
}

func (g *FFmpeg) Open() error {
	stream, err := gocv.OpenVideoCaptureWithAPI(g.device, gocv.VideoCaptureFFmpeg)
	if err != nil {
		return err
	}

	g.stream = stream
	return nil
}

func (g *FFmpeg) Close() error {
	return g.stream.Close()
}

func (g *FFmpeg) Read() (*cv.Frame, error) {
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
