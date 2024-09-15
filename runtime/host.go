package runtime

import (
	"context"
	"log"
	"maps"

	"github.com/orsinium-labs/wypes"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/engine"
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

var guestModules = []api.Module{}

func RegisterGuestModule(ctx context.Context, r wazero.Runtime, module []byte) {
	mod, err := r.Instantiate(ctx, module)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	guestModules = append(guestModules, mod)
}

const process = "process"

func PerformProcessing(ctx context.Context, r wazero.Runtime, frm engine.Frame) engine.Frame {
	var frames []wypes.UInt32

	in := frm.ID
	for _, mod := range guestModules {
		frames = append(frames, wypes.UInt32(in))

		fn := mod.ExportedFunction(process)
		if fn == nil {
			log.Panicf("failed to find function %s", process)
		}

		out, err := fn.Call(ctx, api.EncodeU32(in.Unwrap()))
		if err != nil {
			log.Panicf("failed to call function %s: %v", process, err)
		}
		if len(out) != 1 {
			log.Panicf("expected 1 return value, got %d", len(out))
		}

		in = wypes.UInt32(api.DecodeU32(out[0]))
	}
	out := in

	// close up all the frames except the last one
	for i := 0; i < len(frames)-2; i++ {
		frm := engine.FrameCache[frames[i]]
		frm.Close()
		delete(engine.FrameCache, frames[i])
	}

	return engine.FrameCache[out]
}
