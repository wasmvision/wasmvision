package runtime

import (
	"context"
	"log"
	"maps"

	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/frame"

	"github.com/orsinium-labs/wypes"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Interpreter is a WebAssembly interpreter that can load and run guest modules.
type Interpreter struct {
	r            wazero.Runtime
	guestModules []api.Module
}

// New creates a new Interpreter.
func New(ctx context.Context) Interpreter {
	r := wazero.NewRuntime(ctx)

	modules := hostModules()
	if err := modules.DefineWazero(r, nil); err != nil {
		log.Panicf("error define host functions: %v\n", err)
	}

	return Interpreter{
		r:            r,
		guestModules: []api.Module{},
	}
}

func hostModules() wypes.Modules {
	modules := wypes.Modules{
		"hosted": wypes.Module{
			"println": wypes.H1(hostPrintln),
		},
	}
	maps.Copy(modules, cv.MatModules())
	maps.Copy(modules, cv.ImgprocModules())
	maps.Copy(modules, cv.NetModules())

	return modules
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}

// Close closes the interpreter.
func (intp *Interpreter) Close(ctx context.Context) {
	intp.r.Close(ctx)
}

// RegisterGuestModule registers a guest module with the interpreter.
func (intp *Interpreter) RegisterGuestModule(ctx context.Context, module []byte) error {
	mod, err := intp.r.Instantiate(ctx, module)
	if err != nil {
		return err
	}

	intp.guestModules = append(intp.guestModules, mod)
	return nil
}

const process = "process"

// PerformProcessing performs processing on a frame.
func (intp *Interpreter) PerformProcessing(ctx context.Context, frm frame.Frame) frame.Frame {
	var frames []wypes.UInt32

	in := frm.ID
	for _, mod := range intp.guestModules {
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
