package datastore

import (
	"os"
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

func TestProcessorsBoltDB(t *testing.T) {
	t.Run("set/get", func(t *testing.T) {

		os.Setenv("WASMVISION_STORAGE_BOLTDB_FILENAME", "test-bolt.db")
		defer func() {
			os.Unsetenv("WASMVISION_STORAGE_BOLTDB_FILENAME")
			os.Remove("test-bolt.db")
		}()

		s := NewProcessors(storage.NewBoltDBStorage())

		if s.Error() != nil {
			t.Fatal(s.Error())
		}

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
