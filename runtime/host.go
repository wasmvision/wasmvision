package runtime

import (
	"context"
	"log"
	"maps"
	"os"

	"github.com/wasmvision/wasmvision/config"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/guest"

	"github.com/orsinium-labs/wypes"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Interpreter is a WebAssembly interpreter that can load and run guest modules.
type Interpreter struct {
	r               wazero.Runtime
	Refs            *MapRefs
	guestModules    []guest.Module
	Config          InterpreterConfig
	ModuleContext   *cv.Context
	ProcessorConfig *config.Store
}

type InterpreterConfig struct {
	ProcessorsDir string
	ModelsDir     string
	Logging       bool
}

// New creates a new Interpreter.
func New(ctx context.Context, conf InterpreterConfig) Interpreter {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	configStore := config.NewStore()

	// TODO: populate configStore from config file
	configStore.Set("default", "value")

	cctx := cv.Context{
		ModelsDir: conf.ModelsDir,
		Logging:   conf.Logging,
		Config:    configStore,
	}

	modules := hostModules(&cctx)
	refs := NewMapRefs()
	if err := modules.DefineWazero(r, refs); err != nil {
		log.Panicf("error define host functions: %v\n", err)
	}

	return Interpreter{
		r:               r,
		Refs:            refs,
		guestModules:    []guest.Module{},
		Config:          conf,
		ModuleContext:   &cctx,
		ProcessorConfig: configStore,
	}
}

func hostModules(cctx *cv.Context) wypes.Modules {
	modules := hostedModules(cctx)

	maps.Copy(modules, cv.MatModules(cctx))
	maps.Copy(modules, cv.ImgprocModules(cctx))
	maps.Copy(modules, cv.NetModules(cctx))
	maps.Copy(modules, cv.ObjDetectModules(cctx))

	return modules
}

// Close closes the interpreter.
func (intp *Interpreter) Close(ctx context.Context) {
	intp.r.Close(ctx)
}

func (intp *Interpreter) LoadProcessors(ctx context.Context, processors []string) error {
	for _, p := range processors {
		if guest.ProcessorWellKnown(p) {
			if !guest.ProcessorExists(guest.ProcessorFilename(p, intp.Config.ProcessorsDir)) {
				log.Printf("Downloading processor %s to %s...\n", p, intp.Config.ProcessorsDir)

				if err := guest.DownloadProcessor(p, intp.Config.ProcessorsDir); err != nil {
					return err
				}
			}
		}

		fn := guest.ProcessorFilename(p, intp.Config.ProcessorsDir)

		module, err := os.ReadFile(fn)
		if err != nil {
			return err
		}

		if intp.Config.Logging {
			log.Printf("Loading wasmCV guest module %s...\n", p)
		}

		if err := intp.RegisterGuestModule(ctx, p, module); err != nil {
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
func (intp *Interpreter) RegisterGuestModule(ctx context.Context, name string, module []byte) error {
	mod, err := intp.r.InstantiateWithConfig(ctx, module, wazero.NewModuleConfig().WithName(name).WithStartFunctions("_initialize"))
	if err != nil {
		return err
	}

	// after this we know the ReturnDataPtr for this module
	intp.guestModules = append(intp.guestModules, guest.NewModule(ctx, mod))
	return nil
}

// Process is the exported name of the function in a wasmCV guest module that processes a frame.
const process = "process"

// Process performs processing on a frame.
func (intp *Interpreter) Process(ctx context.Context, frm *cv.Frame) *cv.Frame {
	var frames []wypes.UInt32

	in := frm.ID
	for _, mod := range intp.guestModules {
		frames = append(frames, wypes.UInt32(in))

		fn := mod.ExportedFunction(process)
		if fn == nil {
			log.Panicf("failed to find function %s", process)
		}

		intp.ModuleContext.ReturnDataPtr = mod.ReturnDataPtr

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
	for i := 0; i < len(frames)-1; i++ {
		f, ok := intp.Refs.Get(frames[i].Unwrap(), &cv.Frame{})
		if !ok {
			continue
		}

		fc := f.(*cv.Frame)
		fc.Close()

		intp.Refs.Drop(frames[i].Unwrap())
	}

	f, _ := intp.Refs.Get(out.Unwrap(), &cv.Frame{})

	last := f.(*cv.Frame)

	return last
}
