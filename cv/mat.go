package cv

import (
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func MatModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[constructor]mat":                   wypes.H4(matNewWithSizeFunc(ctx)),
			"[resource-drop]mat":                 wypes.H2(matCloseFunc(ctx)),
			"[static]mat.new-mat":                wypes.H1(matNewFunc(ctx)),
			"[static]mat.new-with-size":          wypes.H4(matNewWithSizeFunc(ctx)),
			"[method]mat.close":                  wypes.H2(matCloseFunc(ctx)),
			"[method]mat.clone":                  wypes.H2(matCloneFunc(ctx)),
			"[method]mat.copy-to":                wypes.H3(matCopyToFunc(ctx)),
			"[method]mat.convert-to":             wypes.H4(matConvertToFunc(ctx)),
			"[method]mat.convert-to-with-params": wypes.H6(matConvertToWithParamsFunc(ctx)),
			"[method]mat.cols":                   wypes.H2(matColsFunc(ctx)),
			"[method]mat.rows":                   wypes.H2(matRowsFunc(ctx)),
			"[method]mat.mattype":                wypes.H2(matTypeFunc(ctx)),
			"[method]mat.empty":                  wypes.H2(matEmptyFunc(ctx)),
			"[method]mat.size":                   wypes.H3(matSizeFunc(ctx)),
			"[method]mat.region":                 wypes.H3(matRegionFunc(ctx)),
			"[method]mat.reshape":                wypes.H5(matReshapeFunc(ctx)),
			"[method]mat.get-float-at":           wypes.H4(matGetFloatAtFunc(ctx)),
			"[method]mat.set-float-at":           wypes.H5(matSetFloatAtFunc(ctx)),
			"[method]mat.get-uchar-at":           wypes.H4(matGetUcharAtFunc(ctx)),
			"[method]mat.set-uchar-at":           wypes.H5(matSetUcharAtFunc(ctx)),
			"[method]mat.get-vecb-at":            wypes.H5(matGetVecbAtFunc(ctx)),
			"[method]mat.add-float":              wypes.H3(matAddFloatFunc(ctx)),
			"[method]mat.subtract-float":         wypes.H3(matSubtractFloatFunc(ctx)),
			"[method]mat.multiply-float":         wypes.H3(matMultiplyFloatFunc(ctx)),
			"[method]mat.divide-float":           wypes.H3(matDivideFloatFunc(ctx)),
			"[method]mat.col-range":              wypes.H5(matColRangeFunc(ctx)),
			"[method]mat.row-range":              wypes.H5(matRowRangeFunc(ctx)),
			"[method]mat.min-max-loc":            wypes.H3(matMinMaxLocFunc(ctx)),
			"[static]mat.ones":                   wypes.H5(matOnesFunc(ctx)),
			"[static]mat.zeros":                  wypes.H5(matZerosFunc(ctx)),
		},
	}
}

func matNewFunc(ctx *Context) func(*wypes.Store) wypes.HostRef[*Frame] {
	return func(s *wypes.Store) wypes.HostRef[*Frame] {
		f := NewEmptyFrame()

		v := wypes.HostRef[*Frame]{Raw: f}

		return v
	}
}

func matNewWithSizeFunc(ctx *Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, rows, cols, matType wypes.UInt32) wypes.HostRef[*Frame] {
		f := NewFrame(gocv.NewMatWithSize(int(rows), int(cols), gocv.MatType(matType)))

		v := wypes.HostRef[*Frame]{Raw: f}

		return v
	}
}

func matCloseFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.Void {
		f := ref.Raw
		f.Close()

		return wypes.Void{}
	}
}

func matColRangeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], start wypes.UInt32, end wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		frm := NewFrame(mat.ColRange(int(start.Unwrap()), int(end.Unwrap())))
		handleFrameReturn(ctx, s, frm, result)
		return wypes.Void{}
	}
}

func matColsFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Cols())
	}
}

func matCloneFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.HostRef[*Frame] {
		f := ref.Raw
		mat := f.Image

		v := wypes.HostRef[*Frame]{Raw: NewFrame(mat.Clone())}
		return v
	}
}

func matCopyToFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.HostRef[*Frame]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], dst wypes.HostRef[*Frame]) wypes.Void {
		f := ref.Raw
		srcMat := f.Image

		dstF := dst.Raw
		dstMat := dstF.Image

		srcMat.CopyTo(&dstMat)

		return wypes.Void{}
	}
}

func matConvertToFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], mattype wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		dst := NewEmptyFrame()

		if err := mat.ConvertTo(&dst.Image, gocv.MatType(mattype.Unwrap())); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func matConvertToWithParamsFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.Float32, wypes.Float32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], mattype wypes.UInt32, alpha wypes.Float32, beta wypes.Float32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		dst := NewEmptyFrame()

		if err := mat.ConvertToWithParams(&dst.Image, gocv.MatType(mattype.Unwrap()), alpha.Unwrap(), beta.Unwrap()); err != nil {
			handleFrameError(ctx, s, dst, result, err)
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func matRegionFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Rect) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], rect Rect) wypes.HostRef[*Frame] {
		f := ref.Raw
		mat := f.Image

		r := rect.Unwrap()
		v := wypes.HostRef[*Frame]{Raw: NewFrame(mat.Region(r))}
		return v
	}
}

func matRowRangeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], start wypes.UInt32, end wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		frm := NewFrame(mat.RowRange(int(start.Unwrap()), int(end.Unwrap())))
		handleFrameReturn(ctx, s, frm, result)
		return wypes.Void{}
	}
}

func matRowsFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Rows())
	}
}

func matTypeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Type())
	}
}

func matEmptyFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame]) wypes.Bool {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame]) wypes.Bool {
		f := ref.Raw
		mat := f.Image

		return wypes.Bool(mat.Empty())
	}
}

func matSizeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.ReturnedList[wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], list wypes.ReturnedList[wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image
		dims := mat.Size()

		result := make([]wypes.UInt32, len(dims))
		for i, dim := range dims {
			result[i] = wypes.UInt32(dim)
		}

		list.Raw = result
		list.DataPtr = ctx.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func matReshapeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], channels wypes.UInt32, rows wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		d := mat.Reshape(int(channels), int(rows))
		dst := NewFrame(d)

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func matGetFloatAtFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.Float32 {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32) wypes.Float32 {
		f := ref.Raw
		mat := f.Image

		return wypes.Float32(mat.GetFloatAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetFloatAtFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Float32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.SetFloatAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetUcharAtFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.UInt8 {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32) wypes.UInt8 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt8(mat.GetUCharAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetUcharAtFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.UInt8) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.UInt8) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.SetUCharAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetVecbAtFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.ReturnedList[wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.ReturnedList[wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image
		data := mat.GetVecbAt(int(row.Unwrap()), int(col.Unwrap()))

		result := make([]wypes.UInt32, len(data))
		for i, dim := range data {
			result[i] = wypes.UInt32(dim)
		}

		v.Raw = result
		v.DataPtr = ctx.ReturnDataPtr
		v.Lower(s)

		return wypes.Void{}
	}
}

func matAddFloatFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.AddFloat(v.Unwrap())

		return wypes.Void{}
	}
}

func matSubtractFloatFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.SubtractFloat(v.Unwrap())

		return wypes.Void{}
	}
}

func matMultiplyFloatFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.MultiplyFloat(v.Unwrap())

		return wypes.Void{}
	}
}

func matDivideFloatFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.DivideFloat(v.Unwrap())

		return wypes.Void{}
	}
}

func matMinMaxLocFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Result[MixMaxLocResult, MixMaxLocResult, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], result wypes.Result[MixMaxLocResult, MixMaxLocResult, wypes.UInt32]) wypes.Void {
		f := ref.Raw
		mat := f.Image

		minVal, maxVal, minLoc, maxLoc := gocv.MinMaxLoc(mat)

		r := MixMaxLocResult{
			MinVal: wypes.Float32(minVal),
			MaxVal: wypes.Float32(maxVal),
			MinLoc: Size{X: wypes.Int32(minLoc.X), Y: wypes.Int32(minLoc.Y)},
			MaxLoc: Size{X: wypes.Int32(maxLoc.X), Y: wypes.Int32(maxLoc.Y)},
		}
		result.OK = r
		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		return wypes.Void{}
	}
}

func matOnesFunc(ctx *Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, rows, cols, matType wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		frm := NewFrame(gocv.Ones(int(rows), int(cols), gocv.MatType(matType)))
		handleFrameReturn(ctx, s, frm, result)

		return wypes.Void{}
	}
}

func matZerosFunc(ctx *Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, rows, cols, matType wypes.UInt32, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		frm := NewFrame(gocv.Zeros(int(rows), int(cols), gocv.MatType(matType)))
		handleFrameReturn(ctx, s, frm, result)

		return wypes.Void{}
	}
}
