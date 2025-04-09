package cv

import (
	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func ImgprocModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/cv": wypes.Module{
			"adaptive-threshold": wypes.H8(adaptiveThresholdFunc(ctx)),
			"blur":               wypes.H4(blurFunc(ctx)),
			"box-filter":         wypes.H5(boxFilterFunc(ctx)),
			"add":                wypes.H4(addFunc(ctx)),
			"subtract":           wypes.H4(subtractFunc(ctx)),
			"multiply":           wypes.H4(multiplyFunc(ctx)),
			"divide":             wypes.H4(divideFunc(ctx)),
			"exp":                wypes.H3(expFunc(ctx)),
			"gaussian-blur":      wypes.H7(gaussianBlurFunc(ctx)),
			"normalize":          wypes.H6(normalizeFunc(ctx)),
			"threshold":          wypes.H6(thresholdFunc(ctx)),
			"resize":             wypes.H7(resizeFunc(ctx)),
			"transpose-ND":       wypes.H4(transposeNDFunc(ctx)),
			"estimate-affine2d":  wypes.H4(estimateAffine2dFunc(ctx)),
			"warp-affine":        wypes.H5(warpAffineFunc(ctx)),
			"put-text":           wypes.H9(putTextFunc(ctx)),
			"rectangle":          wypes.H6(rectangleFunc(ctx)),
			"circle":             wypes.H7(circleFunc(ctx)),
		},
	}
}

func adaptiveThresholdFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], maxValue wypes.Float32, adaptiveThresholdType0 wypes.UInt32, thresholdType0 wypes.UInt32, blockSize0 wypes.UInt32, c0 wypes.Float32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		frm := ref.Raw
		dst := NewEmptyFrame()

		if err := gocv.AdaptiveThreshold(frm.Image, &dst.Image, float32(maxValue), gocv.AdaptiveThresholdType(adaptiveThresholdType0), gocv.ThresholdType(thresholdType0), int(blockSize0), float32(c0)); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func blurFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Size, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], sz Size, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.Blur(src, &dst.Image, sz.Unwrap()); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func boxFilterFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, Size, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], depth0 wypes.UInt32, sz Size, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.BoxFilter(src, &dst.Image, int(depth0), sz.Unwrap()); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func addFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, src1 wypes.HostRef[*Frame], src2 wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		dst := NewEmptyFrame()

		if err := gocv.Add(src1.Raw.Image, src2.Raw.Image, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func subtractFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, src1 wypes.HostRef[*Frame], src2 wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		dst := NewEmptyFrame()

		if err := gocv.Subtract(src1.Raw.Image, src2.Raw.Image, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func multiplyFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, src1 wypes.HostRef[*Frame], src2 wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		dst := NewEmptyFrame()

		if err := gocv.Multiply(src1.Raw.Image, src2.Raw.Image, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func divideFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, src1 wypes.HostRef[*Frame], src2 wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		dst := NewEmptyFrame()

		if err := gocv.Divide(src1.Raw.Image, src2.Raw.Image, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func expFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.Exp(src, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func gaussianBlurFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Size, wypes.Float32, wypes.Float32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], sz Size, sigmaX0 wypes.Float32, sigmaY0 wypes.Float32, border0 wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.GaussianBlur(src, &dst.Image, sz.Unwrap(), float64(sigmaX0), float64(sigmaY0), gocv.BorderType(border0)); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func normalizeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.Float32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], alpha wypes.Float32, beta wypes.Float32, normtype wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.Normalize(src, &dst.Image, float64(alpha), float64(beta), gocv.NormType(normtype)); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func thresholdFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.Float32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], thresh wypes.Float32, maxValue wypes.Float32, thresholdType0 wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		gocv.Threshold(src, &dst.Image, float32(thresh), float32(maxValue), gocv.ThresholdType(thresholdType0))

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func estimateAffine2dFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, from wypes.HostRef[*Frame], to wypes.HostRef[*Frame], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := gocv.NewPoint2fVectorFromMat(from.Raw.Image)
		t := gocv.NewPoint2fVectorFromMat(to.Raw.Image)

		handleFrameReturn(ctx, s, NewFrame(gocv.EstimateAffine2D(f, t)), result)
		return wypes.Void{}
	}
}

func warpAffineFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame], Size, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], m wypes.HostRef[*Frame], sz Size, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.WarpAffine(src, &dst.Image, m.Raw.Image, sz.Unwrap()); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func transposeNDFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.List[wypes.Int32], wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], order wypes.List[wypes.Int32], result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		o := []int{}
		for _, val := range order.Unwrap() {
			o = append(o, int(val.Unwrap()))
		}

		if err := gocv.TransposeND(src, o, &dst.Image); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func resizeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Size, wypes.Float32, wypes.Float32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], sz Size, fx0 wypes.Float32, fy0 wypes.Float32, interp0 wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image
		dst := NewEmptyFrame()

		if err := gocv.Resize(src, &dst.Image, sz.Unwrap(), float64(fx0.Unwrap()), float64(fy0.Unwrap()), gocv.InterpolationFlags(interp0)); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func putTextFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.String, Size, wypes.UInt32, wypes.Float64, RGBA, wypes.UInt32, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], text wypes.String, org Size, fontFace0 wypes.UInt32, fontScale0 wypes.Float64, c RGBA, thickness0 wypes.UInt32, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image

		if err := gocv.PutText(&src, text.Unwrap(), org.Unwrap(), gocv.HersheyFont(fontFace0), float64(fontScale0.Unwrap()), c.Unwrap(), int(thickness0.Unwrap())); err != nil {
			s.Error = err
			result.IsError = true
			result.Error = 1
			result.Lower(s)
		} else {
			result.IsError = false
			result.OK = 0
			result.Lower(s)
		}

		return wypes.Void{}
	}
}

func rectangleFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Rect, RGBA, wypes.UInt32, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], rect Rect, c RGBA, thickness0 wypes.UInt32, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image

		if err := gocv.Rectangle(&src, rect.Unwrap(), c.Unwrap(), int(thickness0.Unwrap())); err != nil {
			s.Error = err
			result.IsError = true
			result.Error = 1
			result.Lower(s)
		} else {
			result.IsError = false
			result.OK = 0
			result.Lower(s)
		}

		return wypes.Void{}
	}
}

func circleFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Size, wypes.UInt32, RGBA, wypes.UInt32, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], size Size, radius wypes.UInt32, c RGBA, thickness0 wypes.UInt32, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		f := ref.Raw
		src := f.Image

		if err := gocv.Circle(&src, size.Unwrap(), int(radius), c.Unwrap(), int(thickness0.Unwrap())); err != nil {
			s.Error = err
			result.IsError = true
			result.Error = 1
			result.Lower(s)
		} else {
			result.IsError = false
			result.OK = 0
			result.Lower(s)
		}

		return wypes.Void{}
	}
}
