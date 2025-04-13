package datastore

import (
	"testing"
)

func TestFrames(t *testing.T) {
	t.Run("set/get", func(t *testing.T) {
		s := NewFrames()

		err := s.Set(101, "key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get(101, "key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})
}
