package engine

import (
	"maps"

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
			"[method]mat.cols": wypes.H1(matColsFunc),
			"[method]mat.rows": wypes.H1(matRowsFunc),
			"[method]mat.type": wypes.H1(matTypeFunc),
		},
		"wasm:cv/cv": wypes.Module{
			"gaussian-blur": wypes.H6(gaussianBlurFunc),
		},
	}
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}
