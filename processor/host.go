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
	maps.Copy(modules, cv.MatModules())
	maps.Copy(modules, cv.ImgprocModules())

	return modules
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}
