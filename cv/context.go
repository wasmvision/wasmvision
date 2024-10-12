package cv

import "github.com/wasmvision/wasmvision/config"

// Context is the configuration for the cv package used when each call is made
// from the guest module.
type Context struct {
	ReturnDataPtr uint32
	ModelsDir     string
	Logging       bool
	Config        *config.Store
}
