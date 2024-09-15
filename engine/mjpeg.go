package engine

import (
	"log"
	"net/http"
	"time"

	"github.com/hybridgroup/mjpeg"
	"github.com/wasmvision/wasmvision/frame"
	"gocv.io/x/gocv"
)

// MJPEGStream represents a Motion JPEG stream used for video streaming display of whatever frames
// are being processed by wasmVision.
type MJPEGStream struct {
	stream *mjpeg.Stream
	Port   string
}

// NewMJPEGStream creates a new MJPEGStream instance with the given port.
func NewMJPEGStream(port string) MJPEGStream {
	return MJPEGStream{Port: port}
}

// Start starts the MJPEG stream server.
func (s *MJPEGStream) Start() {
	s.stream = mjpeg.NewStream()

	http.Handle("/", s.stream)
	server := &http.Server{
		Addr:         s.Port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

// Publish publishes a frame to the MJPEG stream.
func (s *MJPEGStream) Publish(frm frame.Frame) {
	buf, _ := gocv.IMEncode(".jpg", frm.Image)
	defer buf.Close()

	s.stream.UpdateJPEG(buf.GetBytes())
}
