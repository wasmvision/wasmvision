package guest

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter/v2"
)

type ProcessorFile struct {
	Alias       string
	Filename    string
	URL         string
	Description string
}

var knownProcessors = map[string]ProcessorFile{
	"asciify": {
		Alias:       "asciify",
		Filename:    "asciify.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/asciify.wasm",
		Description: "Asciify processor converts frames to ASCII art",
	},
	"blur": {
		Alias:       "blur",
		Filename:    "blur.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blur.wasm",
		Description: "Blur processor applies a blur to frames (Go version)",
	},
	"blurc": {
		Alias:       "blurc",
		Filename:    "blurc.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blurc.wasm",
		Description: "Blurc processor applies a blur to frames (C version)",
	},
	"blurrs": {
		Alias:       "blurrs",
		Filename:    "blurrs.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blurrs.wasm",
		Description: "Blurrs processor applies a blur to frames (Rust version)",
	},
	"captions": {
		Alias:       "captions",
		Filename:    "captions.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/captions.wasm",
		Description: "Captions processor displays captions to frames",
	},
	"edge-detect": {
		Alias:       "edge-detect",
		Filename:    "edge-detect.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/edge-detect.wasm",
		Description: "Edge-detect processor detects edges in frames using a computer vision model",
	},
	"face-expression": {
		Alias:       "face-expression",
		Filename:    "face-expression.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/face-expression.wasm",
		Description: "Face-expression processor performs Facial Expression Recognition on faces that were previously detected using facedetectyn.wasm",
	},
	"faceblur": {
		Alias:       "faceblur",
		Filename:    "faceblur.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/faceblur.wasm",
		Description: "Faceblur processor blurs faces that were previously detected using facedetectyn.wasm",
	},
	"facedetectyn": {
		Alias:       "facedetectyn",
		Filename:    "facedetectyn.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/facedetectyn.wasm",
		Description: "Face detection processor using Yunet vision model",
	},
	"gaussianblur": {
		Alias:       "gaussianblur",
		Filename:    "gaussianblur.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/gaussianblur.wasm",
		Description: "Gaussian blur processor applies a Gaussian blur to frames",
	},
	"hello": {
		Alias:       "hello",
		Filename:    "hello.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/hello.wasm",
		Description: "Hello processor just prints some information about captures frames",
	},
	"object-detector": {
		Alias:       "object-detector",
		Filename:    "object-detector.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/object-detector.wasm",
		Description: "Object detector processor detects objects in frames using YOLOv8 computer vision model",
	},
	"ollama": {
		Alias:       "ollama",
		Filename:    "ollama.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/ollama.wasm",
		Description: "Ollama processor uses the Ollama API to process frames using Language Vision Models",
	},
	"style-transfer": {
		Alias:       "style-transfer",
		Filename:    "style-transfer.wasm",
		URL:         "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/style-transfer.wasm",
		Description: "Style transfer processor applies fast neural style transfer to frames using one of several models",
	},
}

func KnownProcessors() map[string]ProcessorFile {
	return knownProcessors
}

func DownloadProcessor(processor string, processorDir string) error {
	p, ok := knownProcessors[stripWasmExtension(processor)]
	if !ok {
		return errors.New("not known processor")
	}

	req := &getter.Request{
		Src:     p.URL,
		Dst:     filepath.Join(processorDir, filepath.Base(p.Filename)),
		GetMode: getter.ModeFile,
	}

	client := &getter.Client{}
	if _, err := client.Get(context.Background(), req); err != nil {
		return err
	}

	return nil
}

// ProcessorExists checks if the processor file name with full path exists.
func ProcessorExists(processorFile string) bool {
	if _, err := os.Stat(processorFile); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// ProcessorWellKnown checks if the processor is a well-known processor.
// A well-known processor is one that is listed in the knownProcessors map.
// It returns true if the processor is well-known, false otherwise.
func ProcessorWellKnown(processor string) bool {
	processor = stripWasmExtension(processor)
	if _, ok := knownProcessors[processor]; ok {
		return true
	}

	return false
}

// ProcessorFilename returns the full path to the processor file.
// It checks if the processor is a well-known processor or if it exists in the
// specified processor directory. If it is a well-known processor, it returns
// the path to the processor file in the processor directory. If it is not a
// well-known processor, it returns the processor file name as is.
func ProcessorFilename(processor, processorDir string) string {
	p, known := knownProcessors[stripWasmExtension(processor)]
	switch {
	case known:
		// processor name like `asciify` is a well-known processor
		return filepath.Join(processorDir, p.Filename)

	case ProcessorExists(filepath.Join(processorDir, processor)):
		// processor name like `asciify.wasm` is a file in the processors directory
		return filepath.Join(processorDir, processor)

	default:
		// processor name like `/path/to/asciify.wasm` is not a well-known processor and not a file in the processors directory
		return processor
	}
}

func stripWasmExtension(processor string) string {
	if filepath.Ext(processor) == ".wasm" {
		return processor[:len(processor)-len(".wasm")]
	}
	return processor
}
