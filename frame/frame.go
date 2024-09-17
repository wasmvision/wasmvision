package frame

import (
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
func NewFrame() Frame {
	id := rand.IntN(102400)
	return Frame{
		ID: wypes.UInt32(id),
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
