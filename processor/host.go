package processor

import (
	"maps"

	"github.com/hybridgroup/wasmvision/cv"
	"github.com/orsinium-labs/wypes"
)

func HostModules() wypes.Modules {
	modules := wypes.Modules{
		"hosted": wypes.Module{
			"println": wypes.H1(hostPrintln),
		},
	}
	maps.Copy(modules, hostCVModules())

	return modules
}

func hostCVModules() wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[method]mat.cols": wypes.H1(cv.MatColsFunc),
			"[method]mat.rows": wypes.H1(cv.MatRowsFunc),
			"[method]mat.type": wypes.H1(cv.MatTypeFunc),
		},
		"wasm:cv/cv": wypes.Module{
			"gaussian-blur": wypes.H6(cv.GaussianBlurFunc),
		},
	}
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}
