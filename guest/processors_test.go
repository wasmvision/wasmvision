package guest

import (
	"path/filepath"
	"testing"
)

func TestWellKnownProcessor(t *testing.T) {
	t.Run("well-known processor", func(t *testing.T) {
		if !ProcessorWellKnown("style-transfer") {
			t.Errorf("processor style-transfer not found")
		}

		if !ProcessorWellKnown("style-transfer.wasm") {
			t.Errorf("processor style-transfer.wasm not found")
		}
	})

	t.Run("unknown processor", func(t *testing.T) {
		if ProcessorWellKnown("unknown") {
			t.Errorf("processor unknown found")
		}

		if ProcessorWellKnown("unknown.wasm") {
			t.Errorf("processor unknown.wasm found")
		}
	})
}

func TestProcessorFilename(t *testing.T) {
	t.Run("well-known processor", func(t *testing.T) {
		path := filepath.Join("processors", "style-transfer.wasm")

		fn := ProcessorFilename("style-transfer", "processors")
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
