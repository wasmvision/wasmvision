package cv

import (
	"log/slog"
	"math/rand/v2"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

// Frame is a container for an image frame.
type Frame struct {
	ID    wypes.UInt32
	Image gocv.Mat
}

// NewFrame creates a new Frame.
func NewFrame(img gocv.Mat) *Frame {
	return &Frame{
		ID:    wypes.UInt32(rand.IntN(maxIndex)),
		Image: img,
	}
}

// NewEmptyFrame creates a new Frame with an empty Mat.
func NewEmptyFrame() *Frame {
	return &Frame{
		ID:    wypes.UInt32(rand.IntN(maxIndex)),
		Image: gocv.NewMat(),
	}
}

// SetImage sets the image of the frame.
func (f *Frame) SetImage(img gocv.Mat) {
	f.Image = img
}

// Close closes the frame.
func (f *Frame) Close() {
	f.Image.Close()
}

func (f *Frame) Empty() bool {
	return f.Image.Empty()
}

func handleFrameReturn(ctx *Context, s *wypes.Store, frame *Frame, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) {
	result.IsError = false
	result.OK = wypes.HostRef[*Frame]{Raw: frame}
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}

func handleFrameError(ctx *Context, s *wypes.Store, frame *Frame, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32], err error) {
	if err == nil {
		return
	}

	if frame != nil {
		frame.Close()
	}

	slog.Error("cv frame error", "error", err)
	s.Error = err
	result.IsError = true
	result.Error = 1
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}
