package datastore

import (
	"slices"
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

		if string(val) != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("exists", func(t *testing.T) {
		s := NewProcessors(map[string]map[string]string{
			"proc": map[string]string{
				"key": "value",
			},
		})

		ok := s.Exists("proc")
		if !ok {
			t.Errorf("not found")
		}
	})

	t.Run("getKeys", func(t *testing.T) {
		s := NewProcessors(map[string]map[string]string{
			"proc": map[string]string{
				"key":  "value",
				"key2": "value2",
				"key3": "value3",
			},
		})

		keys, ok := s.GetKeys("proc")
		if !ok {
			t.Errorf("processor not found")
		}

		if len(keys) != 3 {
			t.Errorf("unexpected number of keys: %d", len(keys))
		}

		if !slices.Contains(keys, "key") {
			t.Errorf("key not found")
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

		if string(val) != "value" {
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
