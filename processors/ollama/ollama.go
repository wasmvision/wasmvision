//go:build tinygo

package main

import (
	"time"

	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/datastore"
	"github.com/wasmvision/wasmvision-sdk-go/http"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	hosttime "github.com/wasmvision/wasmvision-sdk-go/time"
	"go.bytecodealliance.org/cm"
	"wasmcv.org/wasm/cv/mat"
)

var (
	lastUpdate time.Time

	ps       datastore.ProcessorStore
	template string
	caption  []byte
)

func init() {
	lastUpdate = time.UnixMicro(int64(hosttime.Now(0)))

	ps = datastore.NewProcessorStore(1)
	caption = make([]byte, 0, 80)
}

//export process
func process(image mat.Mat) mat.Mat {
	loadConfig()

	now := time.UnixMicro(int64(hosttime.Now(0)))
	if now.Sub(lastUpdate) > 10*time.Second {
		logging.Info("Asking for image description...")

		req := []byte(template)
		tmpl := cm.ToList[[]byte](req)

		data := http.PostImage(url+"/api/generate", "application/json", tmpl, "response", uint32(image))
		switch {
		case data.IsErr():
			logging.Error("HTTP error")
		case len(data.OK().Slice()) > 0:
			caption = append(caption[:0], data.OK().Slice()...)
			ps.Set("captions", "caption", cm.ToList(caption))
			logging.Println(string(data.OK().Slice()))
		default:
			logging.Info("No result from ollama")
		}

		lastUpdate = now
	}

	return image
}

var (
	url   string
	model string
)

const defaultURL = "http://localhost:11434"
const defaultModel = "llava"

func loadConfig() {
	if url == "" {
		conf := config.GetConfig("url")
		if conf.IsErr() {
			url = defaultURL
		} else {
			url = *conf.OK()
		}

		logging.Info("Using Ollama server at " + url)
	}

	if model == "" {
		conf := config.GetConfig("model")
		if conf.IsErr() {
			model = defaultModel
		} else {
			model = *conf.OK()
		}

		logging.Info("Using Ollama model " + model)
	}

	if template == "" {
		template = `{
			"model": "` + model + `",
			"prompt":"Describe what is in this picture in 60 characters or less.",
			"stream": false,
			"images": ["%IMAGE%"]
		  }`
	}
}
