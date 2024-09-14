package engine

import (
	"github.com/orsinium-labs/wypes"
)

func matColsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Cols())
}

func matRowsFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Rows())
}

func matTypeFunc(matref wypes.UInt32) wypes.UInt32 {
	f, ok := FrameCache[matref]
	if !ok {
		return wypes.UInt32(0)
	}
	mat := f.Image

	return wypes.UInt32(mat.Type())
}
