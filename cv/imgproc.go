package cv

import (
	"image"
	"image/color"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func ImgprocModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/cv": wypes.Module{
			"adaptive-threshold": wypes.H7(adaptiveThresholdFunc(ctx)),
			"blur":               wypes.H4(blurFunc(ctx)),
			"box-filter":         wypes.H5(boxFilterFunc(ctx)),
			"gaussian-blur":      wypes.H7(gaussianBlurFunc(ctx)),
			"threshold":          wypes.H5(thresholdFunc(ctx)),
			"resize":             wypes.H7(resizeFunc(ctx)),
			"put-text":           wypes.H12(putTextFunc(ctx)),
			"rectangle":          wypes.H8(rectangleFunc(ctx)),
			"circle":             wypes.H9(circleFunc(ctx)),
		},
	}
}

func adaptiveThresholdFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], maxValue wypes.Float32, adaptiveThresholdType0 wypes.UInt32, thresholdType0 wypes.UInt32, blockSize0 wypes.UInt32, c0 wypes.Float32) wypes.HostRef[*Frame] {
		frm := ref.Raw

		dst := NewEmptyFrame()

		gocv.AdaptiveThreshold(frm.Image, &dst.Image, float32(maxValue), gocv.AdaptiveThresholdType(adaptiveThresholdType0), gocv.ThresholdType(thresholdType0), int(blockSize0), float32(c0))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func blurFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], size0 wypes.UInt32, size1 wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		src := f.Image

		dst := NewEmptyFrame()

		gocv.Blur(src, &dst.Image, image.Pt(int(size0), int(size1)))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func boxFilterFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], depth0 wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		src := f.Image

		dst := NewEmptyFrame()

		gocv.BoxFilter(src, &dst.Image, int(depth0), image.Pt(int(size0), int(size1)))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func gaussianBlurFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], size0 wypes.UInt32, size1 wypes.UInt32, sigmaX0 wypes.Float32, sigmaY0 wypes.Float32, border0 wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		src := f.Image

		dst := NewEmptyFrame()

		gocv.GaussianBlur(src, &dst.Image, image.Pt(int(size0), int(size1)), float64(sigmaX0), float64(sigmaY0), gocv.BorderType(border0))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func thresholdFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], thresh wypes.Float32, maxValue wypes.Float32, thresholdType0 wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		src := f.Image

		dst := NewEmptyFrame()

		gocv.Threshold(src, &dst.Image, float32(thresh), float32(maxValue), gocv.ThresholdType(thresholdType0))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func resizeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], size0 wypes.UInt32, size1 wypes.UInt32, fx0 wypes.Float32, fy0 wypes.Float32, interp0 wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		src := f.Image

		dst := NewEmptyFrame()

		gocv.Resize(src, &dst.Image, image.Pt(int(size0.Unwrap()), int(size1.Unwrap())), float64(fx0.Unwrap()), float64(fy0.Unwrap()), gocv.InterpolationFlags(interp0))

		v := wypes.HostRef[*Frame]{Raw: dst}
		return v
	}
}

func putTextFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.String, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float64, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], text wypes.String, org0 wypes.UInt32, org1 wypes.UInt32, fontFace0 wypes.UInt32, fontScale0 wypes.Float64, c0 wypes.UInt32, c1 wypes.UInt32, c2 wypes.UInt32, c3 wypes.UInt32, thickness0 wypes.UInt32) wypes.Void {
		f := ref.Raw
		src := f.Image

		clr := color.RGBA{R: uint8(c0.Unwrap()), G: uint8(c1.Unwrap()), B: uint8(c2.Unwrap()), A: uint8(c3.Unwrap())}
		gocv.PutText(&src, text.Unwrap(), image.Pt(int(org0.Unwrap()), int(org1.Unwrap())), gocv.HersheyFont(fontFace0), float64(fontScale0.Unwrap()), clr, int(thickness0.Unwrap()))

		return wypes.Void{}
	}
}

func rectangleFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Rect, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], rect Rect, c0 wypes.UInt32, c1 wypes.UInt32, c2 wypes.UInt32, c3 wypes.UInt32, thickness0 wypes.UInt32) wypes.Void {
		f := ref.Raw
		src := f.Image

		clr := color.RGBA{R: uint8(c0.Unwrap()), G: uint8(c1.Unwrap()), B: uint8(c2.Unwrap()), A: uint8(c3.Unwrap())}
		gocv.Rectangle(&src, rect.Unwrap(), clr, int(thickness0.Unwrap()))

		return wypes.Void{}
	}
}

func circleFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Size, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], size Size, radius wypes.UInt32, c0 wypes.UInt32, c1 wypes.UInt32, c2 wypes.UInt32, c3 wypes.UInt32, thickness0 wypes.UInt32) wypes.Void {
		f := ref.Raw
		src := f.Image

		clr := color.RGBA{R: uint8(c0.Unwrap()), G: uint8(c1.Unwrap()), B: uint8(c2.Unwrap()), A: uint8(c3.Unwrap())}
		gocv.Circle(&src, size.Unwrap(), int(radius), clr, int(thickness0.Unwrap()))

		return wypes.Void{}
	}
}
