package cv

import (
	"github.com/hybridgroup/wasmvision/engine"

	"github.com/orsinium-labs/wypes"
)

func MatModules() wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[method]mat.cols": wypes.H1(matColsFunc),
			"[method]mat.rows": wypes.H1(matRowsFunc),
			"[method]mat.type": wypes.H1(matTypeFunc),
		},
	}
}

func matColsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Cols())
}

func matRowsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Rows())
}

func matTypeFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Type())
}
