//go:build tinygo

package main

import (
	"time"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmvision/wasmvision-sdk-go/config"
	"github.com/wasmvision/wasmvision-sdk-go/http"
	"github.com/wasmvision/wasmvision-sdk-go/logging"
	hosttime "github.com/wasmvision/wasmvision-sdk-go/time"
	"wasmcv.org/wasm/cv/mat"
)

var lastUpdate time.Time

func init() {
	lastUpdate = time.UnixMicro(int64(hosttime.Now(0)))
}

//export process
func process(image mat.Mat) mat.Mat {
	loadConfig()

	now := time.UnixMicro(int64(hosttime.Now(0)))
	if now.Sub(lastUpdate) > 5*time.Second {
		logging.Info("Asking for image description...")

		template := `{
  "model": "` + model + `",
  "prompt":"What is in this picture?",
  "stream": false,
  "images": ["%IMAGE%"]
}`
		req := []byte(template)
		tmpl := cm.ToList[[]byte](req)

		data := http.PostImage(url+"/api/generate", "application/json", tmpl, "response", uint32(image))
		if data.IsErr() {
			logging.Error("HTTP error: " + data.Err().String())
		} else {
			logging.Println(string(data.OK().Slice()))
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
}

func main() {}
