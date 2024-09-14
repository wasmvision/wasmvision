package capture

import (
	"errors"

	"github.com/hybridgroup/wasmvision/engine"
	"gocv.io/x/gocv"
)

type Webcam struct {
	device string
	webcam *gocv.VideoCapture
}

func NewWebcam(device string) Webcam {
	return Webcam{device: device}
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

func (w *Webcam) Read() (engine.Frame, error) {
	img := gocv.NewMat()
	if ok := w.webcam.Read(&img); !ok {
		return engine.Frame{}, errors.New("failed to read frame")
	}

	frame := engine.NewFrame()
	frame.SetImage(img)

	return frame, nil
}
