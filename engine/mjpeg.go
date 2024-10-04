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
	cache  *frame.Cache
	frames chan frame.Frame
}

// NewMJPEGStream creates a new MJPEGStream instance with the given port.
func NewMJPEGStream(cache *frame.Cache, port string) MJPEGStream {
	return MJPEGStream{
		Port:   port,
		cache:  cache,
		frames: make(chan frame.Frame, framebufferSize),
	}
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

	go s.publishFrames()

	log.Fatal(server.ListenAndServe())
}

// Publish publishes a frame to the MJPEG stream.
func (s *MJPEGStream) Publish(frm frame.Frame) error {
	s.frames <- frm
	return nil
}

func (s *MJPEGStream) publishFrames() {
	for frame := range s.frames {
		buf, err := gocv.IMEncode(".jpg", frame.Image)
		if err != nil {
			log.Printf("error writing frame: %v\n", err)
		}
		defer buf.Close()

		s.stream.UpdateJPEG(buf.GetBytes())

		frame.Close()
		s.cache.Delete(frame.ID)
	}
}
