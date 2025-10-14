//go:build llama

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

func vlmInitFromFileFunc[T *VLM](ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, model wypes.String, projector wypes.String, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) wypes.Void {
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

		vlm := NewVLM(modelName, modelFile, projectorFile)
		if err := vlm.Init(); err != nil {
			return handleVLMError(ctx, s, nil, result, fmt.Errorf("cannot init VLM: %v", err))
		}

		return handleVLMSuccess(ctx, s, vlm, result)
	}
}

func vlmCloseFunc(ctx *cv.Context) func(*wypes.Store, wypes.HostRef[*VLM]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*VLM]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}

func vlmPromptFunc(ctx *cv.Context) func(*wypes.Store, wypes.HostRef[*VLM], wypes.String, wypes.HostRef[*cv.Frame], wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, v wypes.HostRef[*VLM], text wypes.String, mat wypes.HostRef[*cv.Frame], result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
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

		if err := vlm.Tokenize(input, bitmap, output); err != nil {
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

	template string
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

func (m *VLM) Init() error {
	slog.Info(fmt.Sprintf("Loading vision language model %s...", m.TextModelFilename))
	m.TextModel = llama.ModelLoadFromFile(m.TextModelFilename, llama.ModelDefaultParams())

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	slog.Info(fmt.Sprintf("Initialize vision language model %s...", m.TextModelFilename))
	m.ModelContext = llama.InitFromModel(m.TextModel, ctxParams)

	m.template = llama.ModelChatTemplate(m.TextModel, "")

	slog.Info("Loading samplers...")
	m.Sampler = llama.NewSampler(m.TextModel, llama.DefaultSamplers)

	slog.Info(fmt.Sprintf("Loading vision language projector %s...", m.ProjectorModelFilename))
	m.ProjectorContext = mtmd.InitFromFile(m.ProjectorModelFilename, m.TextModel, mtmd.ContextParamsDefault())

	return nil
}

func (m *VLM) ChatTemplate(messages []llama.ChatMessage, add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(m.template, messages, add, buf)
	result := string(buf[:len])

	return result
}

func (m *VLM) Tokenize(input *mtmd.InputText, bitmap mtmd.Bitmap, output mtmd.InputChunks) (err error) {
	if res := mtmd.Tokenize(m.ProjectorContext, output, input, []mtmd.Bitmap{bitmap}); res != 0 {
		err = fmt.Errorf("unable to tokenize: %d", res)
	}
	return
}

func (m *VLM) Results(output mtmd.InputChunks) (string, error) {
	var n llama.Pos
	nBatch := 2048 // default value?

	if res := mtmd.HelperEvalChunks(m.ProjectorContext, m.ModelContext, output, 1, 0, int32(nBatch), true, &n); res != 0 {
		return "", errors.New("unable to evaluate chunks")
	}

	var sz int32 = 1
	batch := llama.BatchInit(1, 0, 1)
	batch.NSeqId = &sz
	batch.NTokens = 1
	seqs := unsafe.SliceData([]llama.SeqId{0})
	batch.SeqId = &seqs

	vocab := llama.ModelGetVocab(m.TextModel)
	results := ""

	for i := 0; i < nBatch; i++ {
		token := llama.SamplerSample(m.Sampler, m.ModelContext, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 128)
		len := llama.TokenToPiece(vocab, token, buf, 0, true)
		results += string(buf[:len])

		batch.Token = &token
		batch.Pos = &n

		llama.Decode(m.ModelContext, batch)
		n++
	}

	m.Clear()

	return results, nil
}

// Clear clears the context memory, except for the BOS.
func (m *VLM) Clear() {
	llama.MemorySeqRm(llama.GetMemory(m.ModelContext), 0, 1, -1)
}

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

func handleVLMSuccess(ctx *cv.Context, s *wypes.Store, model *VLM, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32]) wypes.Void {
	result.IsError = false
	result.OK = wypes.HostRef[*VLM]{Raw: model}
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)

	return wypes.Void{}
}

func handleVLMError(ctx *cv.Context, s *wypes.Store, model *VLM, result wypes.Result[wypes.HostRef[*VLM], wypes.HostRef[*VLM], wypes.UInt32], err error) wypes.Void {
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

func handleVLMPromptError(ctx *cv.Context, s *wypes.Store, model *VLM, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32], err error) wypes.Void {
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
