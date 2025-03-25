package engine

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hybridgroup/mjpeg"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/runtime"
	"gocv.io/x/gocv"
)

// MJPEGStream represents a Motion JPEG stream used for video streaming display of whatever frames
// are being processed by wasmVision.
type MJPEGStream struct {
	stream *mjpeg.Stream
	server *http.Server
	Port   string
	refs   *runtime.MapRefs
	frames chan *cv.Frame
}

// NewMJPEGStream creates a new MJPEGStream instance with the given port.
func NewMJPEGStream(refs *runtime.MapRefs, port string) *MJPEGStream {
	return &MJPEGStream{
		Port:   port,
		refs:   refs,
		frames: make(chan *cv.Frame, framebufferSize),
	}
}

// Start starts the MJPEG stream server.
func (s *MJPEGStream) Start() error {
	s.stream = mjpeg.NewStream()

	http.Handle("/", s.stream)
	s.server = &http.Server{
		Addr:         s.Port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go s.publishFrames()
	go func() {
		slog.Error(fmt.Sprintf("mjpeg streamer exited with result %v", s.server.ListenAndServe()))
	}()

	return nil
}

// Close closes the MJPEG stream.
func (s *MJPEGStream) Close() {
	close(s.frames)
	if s.server != nil {
		s.server.Close()
	}
}

// Publish publishes a frame to the MJPEG stream.
func (s *MJPEGStream) Publish(frm *cv.Frame) error {
	s.frames <- frm
	return nil
}

func (s *MJPEGStream) publishFrames() {
	for frame := range s.frames {
		buf, err := gocv.IMEncode(".jpg", frame.Image)
		if err != nil {
			slog.Error(fmt.Sprintf("error writing frame: %v", err))
			continue
		}

		s.stream.UpdateJPEG(buf.GetBytes())

		buf.Close()
		frame.Close()
		s.refs.Drop(frame.ID.Unwrap())
	}
}
