package engine

import (
	"github.com/wasmvision/wasmvision/frame"

	"github.com/orsinium-labs/wypes"
)

var FrameCache = make(map[wypes.UInt32]frame.Frame)
