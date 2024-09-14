package engine

import (
	"image"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func gaussianBlurFunc(matref wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32, sigmaX0 wypes.Float32, sigmaY0 wypes.Float32, border0 wypes.UInt32) wypes.UInt32 {
	f, ok := FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := NewFrame()
	dst.SetImage(gocv.NewMat())
	FrameCache[dst.ID] = dst

	gocv.GaussianBlur(src, &dst.Image, image.Pt(int(size0), int(size1)), float64(sigmaX0), float64(sigmaY0), gocv.BorderType(border0))

	return wypes.UInt32(dst.ID)
}
