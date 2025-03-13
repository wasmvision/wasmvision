package datastore

import (
	"slices"
	"testing"
)

func TestFrames(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{
			145599: map[string]string{
				"key": "value",
			},
		})

		val, ok := s.Get(145599, "key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("getKeys", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{
			145599: map[string]string{
				"key":  "value",
				"key2": "value2",
				"key3": "value3",
			},
		})

		keys, ok := s.GetKeys(145599)
		if !ok {
			t.Errorf("frame not found")
		}

		if len(keys) != 3 {
			t.Errorf("unexpected number of keys: %d", len(keys))
		}

		if !slices.Contains(keys, "key") {
			t.Errorf("key not found")
		}
	})

	t.Run("getKeys non exiting frame", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{
			145599: map[string]string{
				"key":  "value",
				"key2": "value2",
				"key3": "value3",
			},
		})

		_, ok := s.GetKeys(99)
		if ok {
			t.Errorf("frame should not be found")
		}
	})

	t.Run("set", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{})

		err := s.Set(145599, "key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get(145599, "key")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("set multiple keys", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{})

		err := s.Set(145555, "key-1", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = s.Set(145555, "key-2", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = s.Set(145555, "key-3", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get(145555, "key-1")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}

		keys, ok := s.GetKeys(145555)
		if !ok {
			t.Errorf("frame not found")
		}

		if len(keys) != 3 {
			t.Errorf("unexpected number of keys: %d", len(keys))
		}

		if !slices.Contains(keys, "key-2") {
			t.Errorf("key not found")
		}

	})

	t.Run("set multiple keys/frames", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{})

		err := s.Set(1, "key-1", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = s.Set(2, "key-2", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = s.Set(3, "key-3", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		val, ok := s.Get(1, "key-1")
		if !ok {
			t.Errorf("key not found")
		}

		if val != "value" {
			t.Errorf("unexpected value: %s", val)
		}
	})

	t.Run("delete", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{
			1: map[string]string{
				"key": "value",
			},
		})

		s.Delete(1, "key")

		_, ok := s.Get(1, "key")
		if ok {
			t.Errorf("key not deleted")
		}
	})

	t.Run("delete all", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{
			12345: map[string]string{
				"key":  "value",
				"key2": "value2",
				"key3": "value3",
			},
		})

		s.DeleteAll(12345)

		_, ok := s.Get(12345, "key")
		if ok {
			t.Errorf("key not deleted")
		}
	})

	t.Run("exists", func(t *testing.T) {
		s := NewFrames(map[int]map[string]string{})

		err := s.Set(145599, "key", "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		ok := s.Exists(145599)
		if !ok {
			t.Errorf("key not found")
		}
	})
}
