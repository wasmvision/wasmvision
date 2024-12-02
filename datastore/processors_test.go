package datastore

import (
	"testing"
)

func TestProcessors(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		s := NewProcessors(map[string]map[string]string{
			"proc": map[string]string{
				"key": "value",
			},
		})

		val, ok := s.Get("proc", "key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("set", func(t *testing.T) {
		s := NewProcessors(map[string]map[string]string{})

		err := s.Set("proc", "key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get("proc", "key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("delete", func(t *testing.T) {
		s := NewProcessors(map[string]map[string]string{
			"proc": map[string]string{
				"key": "value",
			},
		})

		s.Delete("proc", "key")

		_, ok := s.Get("proc", "key")
		if ok {
			t.Errorf("key not deleted")
		}
	})
}
