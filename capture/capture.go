package capture

import "github.com/wasmvision/wasmvision/frame"

// Capture is the interface that wraps the basic methods for capturing frames.
type Capture interface {
	Open() error
	Close() error
	Read() (frame.Frame, error)
}
