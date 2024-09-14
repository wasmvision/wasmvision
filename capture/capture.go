package capture

import "github.com/hybridgroup/wasmvision/engine"

// Capture is the interface that wraps the basic methods for capturing frames.
type Capture interface {
	Open() error
	Close() error
	Read() (engine.Frame, error)
}
