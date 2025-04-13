package datastore

import (
	"testing"

	"github.com/wasmvision/wasmvision/datastore/storage"
)

func TestProcessors(t *testing.T) {
	t.Run("set/get", func(t *testing.T) {
		s := NewProcessors(storage.NewMemStorage[string]())

		err := s.Set("proc", "key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get("proc", "key")
		if !ok {
			t.Errorf("key not found")
		}

		if string(val) != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})
}
