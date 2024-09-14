package capture

import (
	"errors"

	"github.com/hybridgroup/wasmvision/engine"
	"gocv.io/x/gocv"
)

type WebCam struct {
	device string
	webcam *gocv.VideoCapture
}

func NewWebCam(device string) WebCam {
	return WebCam{device: device}
}

func (w *WebCam) Open() error {
	webcam, err := gocv.OpenVideoCapture(w.device)
	if err != nil {
		return err
	}

	w.webcam = webcam
	return nil
}

func (w *WebCam) Close() error {
	return w.webcam.Close()
}

func (w *WebCam) Read() (engine.Frame, error) {
	img := gocv.NewMat()
	if ok := w.webcam.Read(&img); !ok {
		return engine.Frame{}, errors.New("failed to read frame")
	}

	frame := engine.NewFrame()
	frame.SetImage(img)

	return frame, nil
}
