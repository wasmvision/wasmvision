package datastore

import (
	"os"
	"path"
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

		dbFilename := path.Join(os.TempDir(), "test-bolt.db")

		os.Setenv("WASMVISION_STORAGE_BOLTDB_FILENAME", dbFilename)
		defer func() {
			os.Unsetenv("WASMVISION_STORAGE_BOLTDB_FILENAME")
			os.Remove(dbFilename)
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
