package runtime

import (
	"log"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
)

func hostedModules(ctx *cv.Context) wypes.Modules {
	return wypes.Modules{
		"wasmvision:platform/logging": wypes.Module{
			"println": wypes.H1(hostPrintln),
			"log":     wypes.H1(hostLogFunc(ctx)),
		},
	}
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}

func hostLogFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		if ctx.Logging {
			log.Println(msg.Unwrap())
		}
		return wypes.Void{}
	}
}
