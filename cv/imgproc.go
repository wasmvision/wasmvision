package cv

import (
	"image"

	"github.com/wasmvision/wasmvision/engine"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func ImgprocModules() wypes.Modules {
	return wypes.Modules{
		"wasm:cv/cv": wypes.Module{
			"adaptive-threshold": wypes.H6(adaptiveThresholdFunc),
			"blur":               wypes.H3(blurFunc),
			"box-filter":         wypes.H4(boxFilterFunc),
			"gaussian-blur":      wypes.H6(gaussianBlurFunc),
			"threshold":          wypes.H4(thresholdFunc),
		},
	}
}

func adaptiveThresholdFunc(matref wypes.UInt32, maxValue wypes.Float32, adaptiveThresholdType0 wypes.UInt32, thresholdType0 wypes.UInt32, blockSize0 wypes.UInt32, c0 wypes.Float32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := engine.NewFrame()
	dst.SetImage(gocv.NewMat())
	engine.FrameCache[dst.ID] = dst

	gocv.AdaptiveThreshold(src, &dst.Image, float32(maxValue), gocv.AdaptiveThresholdType(adaptiveThresholdType0), gocv.ThresholdType(thresholdType0), int(blockSize0), float32(c0))

	return wypes.UInt32(dst.ID)
}

func blurFunc(matref wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := engine.NewFrame()
	dst.SetImage(gocv.NewMat())
	engine.FrameCache[dst.ID] = dst

	gocv.Blur(src, &dst.Image, image.Pt(int(size0), int(size1)))

	return wypes.UInt32(dst.ID)
}

func boxFilterFunc(matref wypes.UInt32, depth0 wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := engine.NewFrame()
	dst.SetImage(gocv.NewMat())
	engine.FrameCache[dst.ID] = dst

	gocv.BoxFilter(src, &dst.Image, int(depth0), image.Pt(int(size0), int(size1)))

	return wypes.UInt32(dst.ID)
}

func gaussianBlurFunc(matref wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32, sigmaX0 wypes.Float32, sigmaY0 wypes.Float32, border0 wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := engine.NewFrame()
	dst.SetImage(gocv.NewMat())
	engine.FrameCache[dst.ID] = dst

	gocv.GaussianBlur(src, &dst.Image, image.Pt(int(size0), int(size1)), float64(sigmaX0), float64(sigmaY0), gocv.BorderType(border0))

	return wypes.UInt32(dst.ID)
}

func thresholdFunc(matref wypes.UInt32, thresh wypes.Float32, maxValue wypes.Float32, thresholdType0 wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	src := f.Image

	dst := engine.NewFrame()
	dst.SetImage(gocv.NewMat())
	engine.FrameCache[dst.ID] = dst

	gocv.Threshold(src, &dst.Image, float32(thresh), float32(maxValue), gocv.ThresholdType(thresholdType0))

	return wypes.UInt32(dst.ID)
}
