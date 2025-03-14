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
}
