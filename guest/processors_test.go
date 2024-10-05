package guest

import (
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
		fn := ProcessorFilename("candy", "/tmp")
		if fn != "/tmp/candy.wasm" {
			t.Errorf("unexpected filename %s", fn)
		}
	})

	t.Run("not in processors directory", func(t *testing.T) {
		fn := ProcessorFilename("/some/path/to/unknown.wasm", "/tmp")
		if fn != "/some/path/to/unknown.wasm" {
			t.Errorf("unexpected filename %s", fn)
		}
	})
}
