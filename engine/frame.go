package engine

import (
	"math/rand/v2"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

var FrameCache = make(map[wypes.UInt32]Frame)

type Frame struct {
	ID    wypes.UInt32
	Image gocv.Mat
}

func NewFrame() Frame {
	id := rand.IntN(102400)
	return Frame{
		ID: wypes.UInt32(id),
	}
}

func (f *Frame) SetImage(img gocv.Mat) {
	f.Image = img
}

func (f *Frame) Close() {
	f.Image.Close()
}
