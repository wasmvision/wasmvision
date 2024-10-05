package guest

import (
	"path/filepath"
	"testing"
)

func TestWellKnownProcessor(t *testing.T) {
	t.Run("well-known processor", func(t *testing.T) {
		if !ProcessorWellKnown("candy") {
			t.Errorf("processor candy not found")
		}
	})

	t.Run("unknown processor", func(t *testing.T) {
		if ProcessorWellKnown("unknown") {
			t.Errorf("processor unknown found")
		}
	})
}

func TestProcessorFilename(t *testing.T) {
	t.Run("well-known processor", func(t *testing.T) {
		path := filepath.Join("processors", "candy.wasm")

		fn := ProcessorFilename("candy", "processors")
		if fn != path {
			t.Errorf("unexpected filename %s", fn)
		}
	})

	t.Run("not in processors directory", func(t *testing.T) {
		path := filepath.Join("processors", "are", "awesome", "yes.wasm")

		fn := ProcessorFilename(path, "someplace")
		if fn != path {
			t.Errorf("unexpected filename %s", fn)
		}
	})
}
