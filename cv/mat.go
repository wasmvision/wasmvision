package cv

import (
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func MatModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[constructor]mat":          wypes.H4(matNewWithSizeFunc(ctx)),
			"[resource-drop]mat":        wypes.H2(matCloseFunc(ctx)),
			"[static]mat.new-mat":       wypes.H1(matNewFunc(ctx)),
			"[static]mat.new-with-size": wypes.H4(matNewWithSizeFunc(ctx)),
			"[method]mat.close":         wypes.H2(matCloseFunc(ctx)),
			"[method]mat.clone":         wypes.H2(matCloneFunc(ctx)),
			"[method]mat.copy-to":       wypes.H3(matCopyToFunc(ctx)),
			"[method]mat.cols":          wypes.H2(matColsFunc(ctx)),
			"[method]mat.rows":          wypes.H2(matRowsFunc(ctx)),
			"[method]mat.mattype":       wypes.H2(matTypeFunc(ctx)),
			"[method]mat.empty":         wypes.H2(matEmptyFunc(ctx)),
			"[method]mat.size":          wypes.H3(matSizeFunc(ctx)),
			"[method]mat.region":        wypes.H3(matRegionFunc(ctx)),
			"[method]mat.reshape":       wypes.H4(matReshapeFunc(ctx)),
			"[method]mat.get-float-at":  wypes.H4(matGetFloatAtFunc(ctx)),
			"[method]mat.set-float-at":  wypes.H5(matSetFloatAtFunc(ctx)),
			"[method]mat.get-uchar-at":  wypes.H4(matGetUcharAtFunc(ctx)),
			"[method]mat.set-uchar-at":  wypes.H5(matSetUcharAtFunc(ctx)),
			"[method]mat.get-vecb-at":   wypes.H5(matGetVecbAtFunc(ctx)),
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

func matRegionFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], Rect) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], rect Rect) wypes.HostRef[*Frame] {
		f := ref.Raw
		mat := f.Image

		r := rect.Unwrap()
		v := wypes.HostRef[*Frame]{Raw: NewFrame(mat.Region(r))}
		return v
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

func matReshapeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], channels wypes.UInt32, rows wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		mat := f.Image

		d := mat.Reshape(int(channels), int(rows))
		dst := NewFrame(d)

		v := wypes.HostRef[*Frame]{Raw: dst}

		return v
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
