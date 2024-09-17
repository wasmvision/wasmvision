package cv

import (
	"github.com/wasmvision/wasmvision/frame"

	"github.com/orsinium-labs/wypes"
)

func MatModules(cache *frame.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[method]mat.cols":  wypes.H1(matColsFunc(cache)),
			"[method]mat.rows":  wypes.H1(matRowsFunc(cache)),
			"[method]mat.type":  wypes.H1(matTypeFunc(cache)),
			"[method]mat.empty": wypes.H1(matEmptyFunc(cache)),
		},
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
