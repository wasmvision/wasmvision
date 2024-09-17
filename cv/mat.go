package cv

import (
	"github.com/wasmvision/wasmvision/frame"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func MatModules(cache *frame.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[constructor]mat":    wypes.H3(matNewFunc(cache)),
			"[method]mat.close":   wypes.H1(matCloseFunc(cache)),
			"[method]mat.cols":    wypes.H1(matColsFunc(cache)),
			"[method]mat.rows":    wypes.H1(matRowsFunc(cache)),
			"[method]mat.type":    wypes.H1(matTypeFunc(cache)),
			"[method]mat.empty":   wypes.H1(matEmptyFunc(cache)),
			"[method]mat.reshape": wypes.H3(matReshapeFunc(cache)),
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
