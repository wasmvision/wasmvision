package guest

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter/v2"
)

type ProcessorFile struct {
	Alias    string
	Filename string
	URL      string
}

var knownProcessors = map[string]ProcessorFile{
	"asciify": {
		Alias:    "asciify",
		Filename: "asciify.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/asciify.wasm",
	},
	"blur": {
		Alias:    "blur",
		Filename: "blur.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blur.wasm",
	},
	"blurc": {
		Alias:    "blurc",
		Filename: "blurc.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blurc.wasm",
	},
	"blurrs": {
		Alias:    "blurrs",
		Filename: "blurrs.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/blurrs.wasm",
	},
	"captions": {
		Alias:    "captions",
		Filename: "captions.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/captions.wasm",
	},
	"faceblur": {
		Alias:    "faceblur",
		Filename: "faceblur.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/faceblur.wasm",
	},
	"facedetectyn": {
		Alias:    "facedetectyn",
		Filename: "facedetectyn.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/facedetectyn.wasm",
	},
	"gaussianblur": {
		Alias:    "gaussianblur",
		Filename: "gaussianblur.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/gaussianblur.wasm",
	},
	"hello": {
		Alias:    "hello",
		Filename: "hello.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/hello.wasm",
	},
	"ollama": {
		Alias:    "ollama",
		Filename: "ollama.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/ollama.wasm",
	},
	"style-transfer": {
		Alias:    "style-transfer",
		Filename: "style-transfer.wasm",
		URL:      "https://github.com/wasmvision/wasmvision/raw/refs/heads/main/processors/style-transfer.wasm",
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
