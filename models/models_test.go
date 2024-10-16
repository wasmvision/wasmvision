package models

import (
	"path/filepath"
	"testing"
)

func TestWellKnownModels(t *testing.T) {
	t.Run("well-known model", func(t *testing.T) {
		if !ModelWellKnown("candy-9") {
			t.Errorf("model candy-9 not found")
		}
	})

	t.Run("unknown model", func(t *testing.T) {
		if ModelWellKnown("unknown") {
			t.Errorf("model unknown found")
		}
	})
}

func TestModelFilename(t *testing.T) {
	t.Run("well-known model", func(t *testing.T) {
		path := filepath.Join("models", "candy-9.onnx")

		fn := ModelFileName("candy-9", "models")
		if fn != path {
			t.Errorf("unexpected filename %s", fn)
		}
	})

	t.Run("not well-known", func(t *testing.T) {
		path := filepath.Join("models", "are", "awesome")

		fn := ModelFileName("yes.onnx", path)
		if fn != filepath.Join(path, "yes.onnx") {
			t.Errorf("unexpected filename %s", fn)
		}
	})
}
