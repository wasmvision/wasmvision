package guest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	getter "github.com/hashicorp/go-getter/v2"
	"github.com/wasmvision/wasmvision"
)

type ProcessorFile struct {
	Alias       string
	Filename    string
	Description string
}

var knownProcessors = map[string]ProcessorFile{
	"asciify": {
		Alias:       "asciify",
		Filename:    "asciify.wasm",
		Description: "Asciify processor converts frames to ASCII art",
	},
	"blur": {
		Alias:       "blur",
		Filename:    "blur.wasm",
		Description: "Blur processor applies a blur to frames (Go version)",
	},
	"blurc": {
		Alias:       "blurc",
		Filename:    "blurc.wasm",
		Description: "Blurc processor applies a blur to frames (C version)",
	},
	"blurrs": {
		Alias:       "blurrs",
		Filename:    "blurrs.wasm",
		Description: "Blurrs processor applies a blur to frames (Rust version)",
	},
	"captions": {
		Alias:       "captions",
		Filename:    "captions.wasm",
		Description: "Captions processor displays captions to frames",
	},
	"edge-detect": {
		Alias:       "edge-detect",
		Filename:    "edge-detect.wasm",
		Description: "Edge-detect processor detects edges in frames using a computer vision model",
	},
	"face-counter": {
		Alias:       "face-counter",
		Filename:    "face-counter.wasm",
		Description: "Counts the average number of faces previously detected using facedetectyn.wasm, and saves to Processor datastore.",
	},
	"face-expression": {
		Alias:       "face-expression",
		Filename:    "face-expression.wasm",
		Description: "Face-expression processor performs Facial Expression Recognition on faces that were previously detected using facedetectyn.wasm",
	},
	"faceblur": {
		Alias:       "faceblur",
		Filename:    "faceblur.wasm",
		Description: "Faceblur processor blurs faces that were previously detected using facedetectyn.wasm",
	},
	"facedetectyn": {
		Alias:       "facedetectyn",
		Filename:    "facedetectyn.wasm",
		Description: "Face detection processor using Yunet vision model",
	},
	"gaussianblur": {
		Alias:       "gaussianblur",
		Filename:    "gaussianblur.wasm",
		Description: "Gaussian blur processor applies a Gaussian blur to frames",
	},
	"hello": {
		Alias:       "hello",
		Filename:    "hello.wasm",
		Description: "Hello processor just prints some information about captures frames",
	},
	"object-detector": {
		Alias:       "object-detector",
		Filename:    "object-detector.wasm",
		Description: "Object detector processor detects objects in frames using YOLOv8 computer vision model",
	},
	"ollama": {
		Alias:       "ollama",
		Filename:    "ollama.wasm",
		Description: "Ollama processor uses the Ollama API to process frames using Language Vision Models",
	},
	"style-transfer": {
		Alias:       "style-transfer",
		Filename:    "style-transfer.wasm",
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
		Src:     downloadLocation(p.Filename),
		Dst:     filepath.Join(processorDir, filepath.Base(p.Filename)),
		GetMode: getter.ModeFile,
	}

	client := &getter.Client{}
	if _, err := client.Get(context.Background(), req); err != nil {
		return err
	}

	return nil
}

// downloadLocation returns the download location for the processor file.
// It uses the version of the wasmvision package to determine the correct
// branch to download from. If the version contains "dev", it downloads from
// the current dev branch, otherwise it downloads from the tagged branch for that release.
func downloadLocation(processor string) string {
	version := wasmvision.Version()
	if strings.Contains(version, "dev") {
		return fmt.Sprintf("https://github.com/wasmvision/wasmvision/raw/refs/heads/dev/processors/%s", processor)
	}

	return fmt.Sprintf("https://github.com/wasmvision/wasmvision/raw/refs/tags/%s/processors/%s", version, processor)
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
