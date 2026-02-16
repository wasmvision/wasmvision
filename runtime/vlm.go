//go:build yzma

package runtime

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
	"github.com/hybridgroup/yzma/pkg/vlm"
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"
)

func handleLlama(modules wypes.Modules, enable bool, cctx *cv.Context) error {
	if enable {
		if os.Getenv("YZMA_LIB") == "" {
			return errors.New("YZMA_LIB not set")
		}

		err := llamaInit()
		if err != nil {
			return err
		}
		maps.Copy(modules, hostedVLMModules(cctx))
	}

	return nil
}

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
			"[method]model.prompt":         wypes.H5(vlmPromptFunc(ctx)),
		},
	}
}

func vlmInitFromFileFunc[T *vlm.VLM](ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Result[wypes.HostRef[*vlm.VLM], wypes.HostRef[*vlm.VLM], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, model wypes.String, projector wypes.String, result wypes.Result[wypes.HostRef[*vlm.VLM], wypes.HostRef[*vlm.VLM], wypes.UInt32]) wypes.Void {
		// first the text model
		modelName := model.Unwrap()
		modelFile := models.ModelFileName(modelName, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(modelName):
			slog.Info(fmt.Sprintf("Downloading file for vision language model %s...", modelName))

			if err := models.Download(modelName, ctx.ModelsDir); err != nil {
				return handleVLMError(ctx, s, nil, result, err)
			}

		case !models.ModelExists(modelFile):
			return handleVLMError(ctx, s, nil, result, fmt.Errorf("vision language model %s not found", modelName))
		}

		// now the projector
		projectorName := projector.Unwrap()
		projectorFile := models.ModelFileName(projectorName, ctx.ModelsDir)

		switch {
		case !models.ModelExists(projectorName) && models.ModelWellKnown(projectorName):
			slog.Info(fmt.Sprintf("Downloading file for vision language projector %s...", projectorName))

			if err := models.Download(projectorName, ctx.ModelsDir); err != nil {
				return handleVLMError(ctx, s, nil, result, err)
			}

		case !models.ModelExists(projectorFile):
			return handleVLMError(ctx, s, nil, result, fmt.Errorf("vision language projector %s not found", projectorName))
		}

		vlm := vlm.NewVLM(modelFile, projectorFile)
		if err := vlm.Init(); err != nil {
			return handleVLMError(ctx, s, nil, result, fmt.Errorf("cannot init VLM: %v", err))
		}

		return handleVLMSuccess(ctx, s, vlm, result)
	}
}

func vlmCloseFunc(ctx *cv.Context) func(*wypes.Store, wypes.HostRef[*vlm.VLM]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*vlm.VLM]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}

func vlmPromptFunc(ctx *cv.Context) func(*wypes.Store, wypes.HostRef[*vlm.VLM], wypes.String, wypes.HostRef[*cv.Frame], wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, v wypes.HostRef[*vlm.VLM], text wypes.String, mat wypes.HostRef[*cv.Frame], result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		slog.Info(fmt.Sprintf("prompting vlm with: %s\n", text.Unwrap()))

		vlm := v.Unwrap()
		prompt := text.Unwrap()
		newPrompt := prompt + mtmd.DefaultMarker()

		messages := []llama.ChatMessage{llama.NewChatMessage("user", newPrompt)}
		input := mtmd.NewInputText(vlm.ChatTemplate(messages, true), true, true)

		bitmap, err := matToBitmap(mat.Raw.Image)
		if err != nil {
			return handleVLMPromptError(ctx, s, nil, result, fmt.Errorf("cannot convert image: %v", err))
		}
		defer mtmd.BitmapFree(bitmap)

		output := mtmd.InputChunksInit()
		defer mtmd.InputChunksFree(output)

		if err := vlm.Tokenize(input, []mtmd.Bitmap{bitmap}, output); err != nil {
			return handleVLMPromptError(ctx, s, nil, result, fmt.Errorf("cannot obtain VLM results: %v", err))
		}

		results, err := vlm.Results(output)
		if err != nil {
			return handleVLMPromptError(ctx, s, nil, result, fmt.Errorf("cannot obtain VLM results: %v", err))
		}

		result.IsError = false
		result.OK = wypes.Bytes{Raw: []byte(results)}
		result.DataPtr = ctx.ReturnDataPtr

		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("vlmPromptFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

func llamaInit() error {
	slog.Info("Loading llama.cpp...")
	if err := llama.Load(os.Getenv("YZMA_LIB")); err != nil {
		return err
	}
	if err := mtmd.Load(os.Getenv("YZMA_LIB")); err != nil {
		return err
	}

	slog.Info("Initializing llama.cpp...")
	llama.Init()

	return nil
}

func llamaFree() {
	slog.Info("Unloading llama.cpp...")
	llama.BackendFree()
}

type VlmError uint32

const (
	VlmErrorSuccess VlmError = iota
	VlmErrorRequestError
	VlmErrorRuntimeError
)

func matToBitmap(img gocv.Mat) (mtmd.Bitmap, error) {
	rgb := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	defer rgb.Close()

	gocv.CvtColor(img, &rgb, gocv.ColorBGRToRGB)
	ptr, err := rgb.DataPtrUint8()
	if err != nil {
		return mtmd.Bitmap(0), err
	}

	bitmap := mtmd.BitmapInit(uint32(img.Cols()), uint32(img.Rows()), uintptr(unsafe.Pointer(&ptr[0])))
	return bitmap, nil
}

func handleVLMSuccess(ctx *cv.Context, s *wypes.Store, model *vlm.VLM, result wypes.Result[wypes.HostRef[*vlm.VLM], wypes.HostRef[*vlm.VLM], wypes.UInt32]) wypes.Void {
	result.IsError = false
	result.OK = wypes.HostRef[*vlm.VLM]{Raw: model}
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)

	return wypes.Void{}
}

func handleVLMError(ctx *cv.Context, s *wypes.Store, model *vlm.VLM, result wypes.Result[wypes.HostRef[*vlm.VLM], wypes.HostRef[*vlm.VLM], wypes.UInt32], err error) wypes.Void {
	if err == nil {
		return wypes.Void{}
	}

	slog.Error("VLM error", "error", err)
	s.Error = err
	result.IsError = true
	result.Error = wypes.UInt32(VlmErrorRequestError)
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)

	return wypes.Void{}
}

func handleVLMPromptError(ctx *cv.Context, s *wypes.Store, model *vlm.VLM, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32], err error) wypes.Void {
	if err == nil {
		return wypes.Void{}
	}

	slog.Error("VLM error", "error", err)
	s.Error = err
	result.IsError = true
	result.Error = wypes.UInt32(VlmErrorRequestError)
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)

	return wypes.Void{}
}
