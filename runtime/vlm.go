package runtime

import (
	"fmt"
	"log/slog"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/models"
)

// hostedVLMModules returns the modules that the host provides to the guest
// for using Vision Language Models.
// These are all defined in the wasmvision platform sdk.
// See https://github.com/wasmvision/wasmvision-sdk
func hostedVLMModules(ctx *cv.Context) wypes.Modules {
	return wypes.Modules{
		"wasmvision:platform/vlm": wypes.Module{
			"[static]model.init-from-file": wypes.H4(vlmInitFromFileFunc(ctx)),
			"[resource-drop]model":         wypes.H2(vlmCloseFunc(ctx)),
			"[method]model.close":          wypes.H2(vlmCloseFunc(ctx)),
		},
	}
}

func vlmInitFromFileFunc[T *VLM](ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, model wypes.String, projector wypes.String, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) wypes.Void {
		// first the text model
		modelName := model.Unwrap()
		modelFile := models.ModelFileName(modelName, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(modelName):
			slog.Info(fmt.Sprintf("Downloading file for vision language model %s...", modelName))

			if err := models.Download(modelName, ctx.ModelsDir); err != nil {
				handleVLMError(ctx, s, nil, result, err)
				return wypes.Void{}
			}

		case !models.ModelExists(modelFile):
			handleVLMError(ctx, s, nil, result, fmt.Errorf("vision language model %s not found", modelName))
			return wypes.Void{}
		}

		// now the projector
		projectorName := projector.Unwrap()
		projectorFile := models.ModelFileName(projectorName, ctx.ModelsDir)

		switch {
		case !models.ModelExists(projectorName) && models.ModelWellKnown(projectorName):
			slog.Info(fmt.Sprintf("Downloading file for vision language projector %s...", projectorName))

			if err := models.Download(projectorName, ctx.ModelsDir); err != nil {
				handleVLMError(ctx, s, nil, result, err)
				return wypes.Void{}
			}

		case !models.ModelExists(projectorFile):
			handleVLMError(ctx, s, nil, result, fmt.Errorf("vision language projector %s not found", projectorName))
			return wypes.Void{}
		}

		vlm := NewVLM(modelName, modelFile, projectorFile)

		slog.Info(fmt.Sprintf("Loading vision language model %s...", modelName))
		vlm.TextModel = llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())

		ctxParams := llama.ContextDefaultParams()
		ctxParams.NCtx = 4096
		ctxParams.NBatch = 2048

		slog.Info(fmt.Sprintf("Loading vision language projector %s...", projectorName))
		vlm.ModelContext = llama.InitFromModel(vlm.TextModel, ctxParams)

		vlm.Sampler = llama.NewSampler(vlm.TextModel, llama.DefaultSamplers)
		vlm.ProjectorContext = mtmd.InitFromFile(projectorFile, vlm.TextModel, mtmd.ContextParamsDefault())

		handleVLMReturn(ctx, s, vlm, result)
		return wypes.Void{}
	}
}

func vlmCloseFunc(ctx *cv.Context) func(*wypes.Store, wypes.HostRef[*VLM]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*VLM]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}

// VLM is a Vision Language Model (VLM).
type VLM struct {
	ID                     wypes.UInt32
	Name                   string
	TextModelFilename      string
	ProjectorModelFilename string

	TextModel        llama.Model
	Sampler          llama.Sampler
	ModelContext     llama.Context
	ProjectorContext mtmd.Context
}

// NewVLM creates a new VLM.
func NewVLM(name, model, projector string) *VLM {
	return &VLM{
		Name:                   name,
		TextModelFilename:      model,
		ProjectorModelFilename: projector,
	}
}

// Close closes the VLM.
func (m *VLM) Close() {
	if m.ProjectorContext != mtmd.Context(0) {
		mtmd.Free(m.ProjectorContext)

	}

	if m.ModelContext != llama.Context(0) {
		llama.Free(m.ModelContext)

	}
}

func handleVLMReturn(ctx *cv.Context, s *wypes.Store, model *VLM, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) {
	result.IsError = false
	result.OK = wypes.HostRef[*VLM]{Raw: model}
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}

func handleVLMError(ctx *cv.Context, s *wypes.Store, model *VLM, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32], err error) {
	if err == nil {
		return
	}

	slog.Error("VLM error", "error", err)
	s.Error = err
	result.IsError = true
	result.Error = 1
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}
