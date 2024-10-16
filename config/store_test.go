package config

import (
	"testing"
)

func TestStore(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		s := NewStore(map[string]string{
			"key": "value",
		})

		val, ok := s.Get("key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("set", func(t *testing.T) {
		s := NewStore(map[string]string{})

		err := s.Set("key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get("key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("delete", func(t *testing.T) {
		s := NewStore(map[string]string{
			"key": "value",
		})

		s.Delete("key")

		_, ok := s.Get("key")
		if ok {
			t.Errorf("key not deleted")
		}
	})
}
