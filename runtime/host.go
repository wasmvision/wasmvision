package runtime

import (
	"context"
	"log"
	"maps"

	"github.com/hybridgroup/wasmvision/cv"
	"github.com/orsinium-labs/wypes"
	"github.com/tetratelabs/wazero"
)

func New(ctx context.Context) wazero.Runtime {
	r := wazero.NewRuntime(ctx)

	modules := hostModules()
	if err := modules.DefineWazero(r, nil); err != nil {
		log.Panicf("error define host functions: %v\n", err)
	}

	return r
}

func hostModules() wypes.Modules {
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
