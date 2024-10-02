package capture

import (
	"errors"

	"github.com/wasmvision/wasmvision/frame"
)

var ErrClosed = errors.New("capture device closed")

// Capture is the interface that wraps the basic methods for capturing frames.
type Capture interface {
	Open() error
	Close() error
	Read() (frame.Frame, error)
}
