package capture

import (
	"github.com/wasmvision/wasmvision/frame"

	"gocv.io/x/gocv"
)

type Webcam struct {
	device string
	webcam *gocv.VideoCapture
}

func NewWebcam(device string) *Webcam {
	return &Webcam{device: device}
}

func (w *Webcam) Open() error {
	webcam, err := gocv.OpenVideoCapture(w.device)
	if err != nil {
		return err
	}

	w.webcam = webcam
	return nil
}

func (w *Webcam) Close() error {
	return w.webcam.Close()
}

func (w *Webcam) Read() (frame.Frame, error) {
	img := gocv.NewMat()
	if ok := w.webcam.Read(&img); !ok {
		return frame.Frame{}, ErrClosed
	}

	frame := frame.NewFrame()
	frame.SetImage(img)

	return frame, nil
}
