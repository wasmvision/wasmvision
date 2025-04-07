package cv

import (
	"github.com/wasmvision/wasmvision/config"
	"github.com/wasmvision/wasmvision/datastore"
)

// Context is the configuration for the cv package used when each call is made
// from the guest module.
type Context struct {
	ReturnDataPtr  uint32
	ModelsDir      string
	Config         *config.Store
	FrameStore     *datastore.Frames
	ProcessorStore *datastore.Processors
	EnableCUDA     bool
}

func NewContext(modelsDir string, conf *config.Store, enableCUDA bool) *Context {
	return &Context{
		ModelsDir:      modelsDir,
		Config:         conf,
		FrameStore:     datastore.NewFrames(map[int]map[string]string{}),
		ProcessorStore: datastore.NewProcessors(map[string]map[string]string{}),
		EnableCUDA:     enableCUDA,
	}
}
