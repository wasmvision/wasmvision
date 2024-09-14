package cv

import (
	"github.com/hybridgroup/wasmvision/engine"

	"github.com/orsinium-labs/wypes"
)

func MatColsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Cols())
}

func MatRowsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Rows())
}

func MatTypeFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := engine.FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Type())
}
