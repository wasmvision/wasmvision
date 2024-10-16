package engine

import (
	"testing"

	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/runtime"
)

func TestMJPEGStream(t *testing.T) {
	t.Run("new mjpeg stream", func(t *testing.T) {
		refs := runtime.NewMapRefs()
		port := ":8080"

		s := NewMJPEGStream(refs, port)

		if s.Port != port {
			t.Errorf("unexpected port: %s", s.Port)
		}

		if s.refs != refs {
			t.Errorf("unexpected refs")
		}

		if s.frames == nil {
			t.Errorf("unexpected nil frames")
		}
	})
}

func TestMJPEGStreamStart(t *testing.T) {
	t.Run("start mjpeg stream", func(t *testing.T) {
		refs := runtime.NewMapRefs()
		port := ":8080"

		s := NewMJPEGStream(refs, port)

		defer s.Close()

		go s.Start()

		frm := cv.NewEmptyFrame()

		if err := s.Publish(frm); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
