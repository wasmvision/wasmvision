package cv

import (
	"github.com/wasmvision/wasmvision/frame"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func MatModules(cache *frame.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[constructor]mat":         wypes.H3(matNewFunc(cache)),
			"[resource-drop]mat":       wypes.H1(matCloseFunc(cache)),
			"[method]mat.close":        wypes.H1(matCloseFunc(cache)),
			"[method]mat.cols":         wypes.H1(matColsFunc(cache)),
			"[method]mat.rows":         wypes.H1(matRowsFunc(cache)),
			"[method]mat.mattype":      wypes.H1(matTypeFunc(cache)),
			"[method]mat.empty":        wypes.H1(matEmptyFunc(cache)),
			"[method]mat.size":         wypes.H3(matSizeFunc(cache)),
			"[method]mat.reshape":      wypes.H3(matReshapeFunc(cache)),
			"[method]mat.get-float-at": wypes.H3(matGetFloatAtFunc(cache)),
			"[method]mat.set-float-at": wypes.H4(matSetFloatAtFunc(cache)),
			"[method]mat.get-uchar-at": wypes.H3(matGetUcharAtFunc(cache)),
			"[method]mat.set-uchar-at": wypes.H4(matSetUcharAtFunc(cache)),
			"[method]mat.get-vecb-at":  wypes.H5(matGetVecbAtFunc(cache)),
		},
	}
}

func matNewFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.UInt32 {
	return func(rows, cols, matType wypes.UInt32) wypes.UInt32 {
		mat := gocv.NewMatWithSize(int(rows), int(cols), gocv.MatType(matType))

		f := frame.NewFrame()
		f.SetImage(mat)

		cache.Set(f)

		return f.ID
	}
}

func matCloseFunc(cache *frame.Cache) func(wypes.UInt32) wypes.Void {
	return func(ref wypes.UInt32) wypes.Void {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		mat := f.Image
		mat.Close()

		return wypes.Void{}
	}
}

func matColsFunc(cache *frame.Cache) func(wypes.UInt32) wypes.UInt32 {
	return func(ref wypes.UInt32) wypes.UInt32 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt32(0)
		}
		mat := f.Image

		return wypes.UInt32(mat.Cols())
	}
}

func matRowsFunc(cache *frame.Cache) func(wypes.UInt32) wypes.UInt32 {
	return func(ref wypes.UInt32) wypes.UInt32 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt32(0)
		}
		mat := f.Image

		return wypes.UInt32(mat.Rows())
	}
}

func matTypeFunc(cache *frame.Cache) func(wypes.UInt32) wypes.UInt32 {
	return func(ref wypes.UInt32) wypes.UInt32 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt32(0)
		}
		mat := f.Image

		return wypes.UInt32(mat.Type())
	}
}

func matEmptyFunc(cache *frame.Cache) func(wypes.UInt32) wypes.Bool {
	return func(ref wypes.UInt32) wypes.Bool {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Bool(true)
		}
		mat := f.Image

		return wypes.Bool(mat.Empty())
	}
}

func matSizeFunc(cache *frame.Cache) func(wypes.Store, wypes.UInt32, wypes.List[uint32]) wypes.Void {
	return func(s wypes.Store, ref wypes.UInt32, list wypes.List[uint32]) wypes.Void {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		mat := f.Image
		dims := mat.Size()

		result := make([]uint32, len(dims))
		for i, dim := range dims {
			result[i] = uint32(dim)
		}

		list.Raw = result
		list.DataPtr = cache.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func matReshapeFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.UInt32 {
	return func(ref wypes.UInt32, channels wypes.UInt32, rows wypes.UInt32) wypes.UInt32 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt32(0)
		}
		mat := f.Image

		d := mat.Reshape(int(channels), int(rows))
		dst := frame.NewFrame()
		dst.SetImage(d)

		cache.Set(dst)

		return dst.ID
	}
}

func matGetFloatAtFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.Float32 {
	return func(ref wypes.UInt32, row wypes.UInt32, col wypes.UInt32) wypes.Float32 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Float32(0)
		}
		mat := f.Image

		return wypes.Float32(mat.GetFloatAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetFloatAtFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.Float32) wypes.Void {
	return func(ref wypes.UInt32, row wypes.UInt32, col wypes.UInt32, v wypes.Float32) wypes.Void {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		mat := f.Image
		mat.SetFloatAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetUcharAtFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.UInt8 {
	return func(ref wypes.UInt32, row wypes.UInt32, col wypes.UInt32) wypes.UInt8 {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt8(0)
		}
		mat := f.Image

		return wypes.UInt8(mat.GetUCharAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetUcharAtFunc(cache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.UInt8) wypes.Void {
	return func(ref wypes.UInt32, row wypes.UInt32, col wypes.UInt32, v wypes.UInt8) wypes.Void {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		mat := f.Image
		mat.SetUCharAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetVecbAtFunc(cache *frame.Cache) func(wypes.Store, wypes.UInt32, wypes.UInt32, wypes.UInt32, wypes.List[uint8]) wypes.Void {
	return func(s wypes.Store, ref wypes.UInt32, row wypes.UInt32, col wypes.UInt32, v wypes.List[uint8]) wypes.Void {
		f, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		mat := f.Image
		data := mat.GetVecbAt(int(row.Unwrap()), int(col.Unwrap()))

		result := make([]uint8, len(data))
		for i, dim := range data {
			result[i] = uint8(dim)
		}

		v.Raw = result
		v.DataPtr = cache.ReturnDataPtr
		v.Lower(s)

		return wypes.Void{}
	}
}
