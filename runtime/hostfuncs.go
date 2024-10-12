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
		"wasmvision:platform/config": wypes.Module{
			"get-config": wypes.H3(hostGetConfig(ctx)),
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

func hostGetConfig(ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, key wypes.String, result wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
		if s.Error != nil {
			log.Printf("error in store after lift: %v\n", s.Error)
		}
		val, ok := ctx.Config.Get(key.Unwrap())
		if !ok {
			result.IsError = true
			result.Error = wypes.UInt32(1) // no-such-key
		} else {
			result.IsError = false
			result.OK = wypes.String{Raw: val}
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		if s.Error != nil {
			log.Printf("error in store after lower: %v\n", s.Error)
		}

		return wypes.Void{}
	}
}
