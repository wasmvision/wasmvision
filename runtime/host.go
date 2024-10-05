package runtime

import (
	"context"
	"log"
	"maps"
	"os"

	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/frame"
	"github.com/wasmvision/wasmvision/guest"
	"github.com/wasmvision/wasmvision/net"

	"github.com/orsinium-labs/wypes"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Interpreter is a WebAssembly interpreter that can load and run guest modules.
type Interpreter struct {
	r             wazero.Runtime
	guestModules  []guest.Module
	FrameCache    *frame.Cache
	NetCache      *net.Cache
	ProcessorsDir string
	Logging       bool
}

type InterpreterConfig struct {
	ProcessorsDir string
	ModelsDir     string
	Logging       bool
}

// New creates a new Interpreter.
func New(ctx context.Context, config InterpreterConfig) Interpreter {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	cache := frame.NewCache()
	nc := net.NewCache()
	nc.ModelsDir = config.ModelsDir

	modules := hostModules(cache, nc, config.Logging)
	if err := modules.DefineWazero(r, nil); err != nil {
		log.Panicf("error define host functions: %v\n", err)
	}

	return Interpreter{
		r:             r,
		guestModules:  []guest.Module{},
		FrameCache:    cache,
		NetCache:      nc,
		ProcessorsDir: config.ProcessorsDir,
		Logging:       config.Logging,
	}
}

func hostModules(cache *frame.Cache, nc *net.Cache, logging bool) wypes.Modules {
	modules := hostedModules(logging)
	maps.Copy(modules, cv.MatModules(cache))
	maps.Copy(modules, cv.ImgprocModules(cache))
	maps.Copy(modules, cv.NetModules(cache, nc))

	return modules
}

// Close closes the interpreter.
func (intp *Interpreter) Close(ctx context.Context) {
	intp.r.Close(ctx)
}

func (intp *Interpreter) LoadProcessors(ctx context.Context, processors []string) error {
	for _, p := range processors {
		if guest.ProcessorWellKnown(p) {
			if !guest.ProcessorExists(guest.ProcessorFilename(p, intp.ProcessorsDir)) {
				log.Printf("Downloading processor %s to %s...\n", p, intp.ProcessorsDir)

				if err := guest.DownloadProcessor(p, intp.ProcessorsDir); err != nil {
					return err
				}
			}
		}

		fn := guest.ProcessorFilename(p, intp.ProcessorsDir)

		module, err := os.ReadFile(fn)
		if err != nil {
			return err
		}

		if intp.Logging {
			log.Printf("Loading wasmCV guest module %s...\n", p)
		}

		if err := intp.RegisterGuestModule(ctx, module); err != nil {
			return err
		}
	}

	return nil
}

// Processors returns the guest modules registered with the interpreter.
func (intp *Interpreter) Processors() []guest.Module {
	return intp.guestModules
}

// RegisterGuestModule registers a guest module with the interpreter.
func (intp *Interpreter) RegisterGuestModule(ctx context.Context, module []byte) error {
	mod, err := intp.r.InstantiateWithConfig(ctx, module, wazero.NewModuleConfig().WithName("").WithStartFunctions("_initialize"))
	if err != nil {
		println("error instantiate module")
		return err
	}

	intp.guestModules = append(intp.guestModules, guest.NewModule(ctx, mod))
	return nil
}

// Process is the exported name of the function in a wasmCV guest module that processes a frame.
const process = "process"

// Process performs processing on a frame.
func (intp *Interpreter) Process(ctx context.Context, frm frame.Frame) frame.Frame {
	var frames []wypes.UInt32

	in := frm.ID
	for _, mod := range intp.guestModules {
		frames = append(frames, wypes.UInt32(in))

		fn := mod.ExportedFunction(process)
		if fn == nil {
			log.Panicf("failed to find function %s", process)
		}

		intp.FrameCache.ReturnDataPtr = mod.ReturnDataPtr
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
		frm, _ := intp.FrameCache.Get(frames[i])
		frm.Close()

		intp.FrameCache.Delete(frames[i])
	}

	last, _ := intp.FrameCache.Get(out)
	return last
}
