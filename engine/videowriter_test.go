package engine

import (
	"testing"

	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/runtime"
	"gocv.io/x/gocv"
)

func TestVideoWriter(t *testing.T) {
	t.Run("new video writer", func(t *testing.T) {
		refs := runtime.NewMapRefs()
		dest := "output.avi"

		vw := NewVideoWriter(refs, dest)

		if vw.Filename != dest {
			t.Errorf("unexpected filename: %s", vw.Filename)
		}

		if vw.codec != defaultCodec {
			t.Errorf("unexpected codec: %s", vw.codec)
		}

		if vw.fps != defaultFBS {
			t.Errorf("unexpected fps: %f", vw.fps)
		}

		if vw.refs != refs {
			t.Errorf("unexpected refs")
		}

		if vw.frames == nil {
			t.Errorf("unexpected nil frames")
		}
	})

	t.Run("write video frame", func(t *testing.T) {
		refs := runtime.NewMapRefs()

		tmp := t.TempDir()
		dest := tmp + "/output.avi"

		vw := NewVideoWriter(refs, dest)

		defer vw.Close()

		img := gocv.IMRead("../images/wasmvision-logo.png", gocv.IMReadColor)
		frm := cv.NewFrame(img)

		if err := vw.Write(frm); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
