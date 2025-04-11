package capture

import (
	"github.com/wasmvision/wasmvision/cv"

	"gocv.io/x/gocv"
)

type Webcam struct {
	device  string
	webcam  *gocv.VideoCapture
	retries int
}

func NewWebcam(device string) *Webcam {
	return &Webcam{
		device:  device,
		retries: defaultRetries,
	}
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

func (w *Webcam) Read() (*cv.Frame, error) {
	img := gocv.NewMat()
	if ok := w.webcam.Read(&img); !ok {
		w.retries--
		if w.retries == 0 {
			return &cv.Frame{}, ErrClosed
		}

		frame := cv.NewEmptyFrame()

		return frame, nil
	}

	w.retries = defaultRetries

	frame := cv.NewFrame(img)

	return frame, nil
}

func (w *Webcam) Get(property int32) (value float32, err error) {
	return float32(w.webcam.Get(gocv.VideoCaptureProperties(property))), nil
}

func (w *Webcam) Set(property int32, value float32) error {
	w.webcam.Set(gocv.VideoCaptureProperties(property), float64(value))
	return nil
}
