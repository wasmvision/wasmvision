package cv

import (
	"image"
	"image/color"

	"github.com/wasmvision/wasmvision/frame"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func ImgprocModules(config *Config) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/cv": wypes.Module{
			"adaptive-threshold": wypes.H7(adaptiveThresholdFunc(config)),
			"blur":               wypes.H4(blurFunc(config)),
			"box-filter":         wypes.H5(boxFilterFunc(config)),
			"gaussian-blur":      wypes.H7(gaussianBlurFunc(config)),
			"threshold":          wypes.H5(thresholdFunc(config)),
			"resize":             wypes.H7(resizeFunc(config)),
			"put-text":           wypes.H12(putTextFunc(config)),
		},
	}
}

func adaptiveThresholdFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], maxValue wypes.Float32, adaptiveThresholdType0 wypes.UInt32, thresholdType0 wypes.UInt32, blockSize0 wypes.UInt32, c0 wypes.Float32) wypes.HostRef[*frame.Frame] {
		frm := ref.Raw

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.AdaptiveThreshold(frm.Image, &dst.Image, float32(maxValue), gocv.AdaptiveThresholdType(adaptiveThresholdType0), gocv.ThresholdType(thresholdType0), int(blockSize0), float32(c0))

		return v
	}
}

func blurFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.UInt32, wypes.UInt32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], size0 wypes.UInt32, size1 wypes.UInt32) wypes.HostRef[*frame.Frame] {
		f := ref.Raw
		src := f.Image

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.Blur(src, &dst.Image, image.Pt(int(size0), int(size1)))

		return v
	}
}

func boxFilterFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], depth0 wypes.UInt32, size0 wypes.UInt32, size1 wypes.UInt32) wypes.HostRef[*frame.Frame] {
		f := ref.Raw
		src := f.Image

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.BoxFilter(src, &dst.Image, int(depth0), image.Pt(int(size0), int(size1)))

		return v
	}
}

func gaussianBlurFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], size0 wypes.UInt32, size1 wypes.UInt32, sigmaX0 wypes.Float32, sigmaY0 wypes.Float32, border0 wypes.UInt32) wypes.HostRef[*frame.Frame] {
		f := ref.Raw
		src := f.Image

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.GaussianBlur(src, &dst.Image, image.Pt(int(size0), int(size1)), float64(sigmaX0), float64(sigmaY0), gocv.BorderType(border0))

		return v
	}
}

func thresholdFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], thresh wypes.Float32, maxValue wypes.Float32, thresholdType0 wypes.UInt32) wypes.HostRef[*frame.Frame] {
		f := ref.Raw
		src := f.Image

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.Threshold(src, &dst.Image, float32(thresh), float32(maxValue), gocv.ThresholdType(thresholdType0))

		return v
	}
}

func resizeFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.UInt32) wypes.HostRef[*frame.Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], size0 wypes.UInt32, size1 wypes.UInt32, fx0 wypes.Float32, fy0 wypes.Float32, interp0 wypes.UInt32) wypes.HostRef[*frame.Frame] {
		f := ref.Raw
		src := f.Image

		dst := frame.NewEmptyFrame()
		v := wypes.HostRef[*frame.Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		gocv.Resize(src, &dst.Image, image.Pt(int(size0.Unwrap()), int(size1.Unwrap())), float64(fx0.Unwrap()), float64(fy0.Unwrap()), gocv.InterpolationFlags(interp0))

		return v
	}
}

func putTextFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame], wypes.String, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float64, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.Void {
	return func(store wypes.Store, ref wypes.HostRef[*frame.Frame], text wypes.String, org0 wypes.UInt32, org1 wypes.UInt32, fontFace0 wypes.UInt32, fontScale0 wypes.Float64, c0 wypes.UInt32, c1 wypes.UInt32, c2 wypes.UInt32, c3 wypes.UInt32, thickness0 wypes.UInt32) wypes.Void {
		f := ref.Raw
		src := f.Image

		clr := color.RGBA{R: uint8(c0.Unwrap()), G: uint8(c1.Unwrap()), B: uint8(c2.Unwrap()), A: uint8(c3.Unwrap())}
		gocv.PutText(&src, text.Unwrap(), image.Pt(int(org0.Unwrap()), int(org1.Unwrap())), gocv.HersheyFont(fontFace0), float64(fontScale0.Unwrap()), clr, int(thickness0.Unwrap()))

		return wypes.Void{}
	}
}
