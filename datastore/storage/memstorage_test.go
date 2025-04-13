package storage

import (
	"testing"
)

func TestMemStorageSet(t *testing.T) {
	store := NewMemStorage[string]()

	err := store.Set("key1", "subkey1", "value1")
	if err != nil {
		t.Errorf("Set() returned an error: %v", err)
	}

	val, ok := store.Get("key1", "subkey1")
	if !ok || val != "value1" {
		t.Errorf("Get() = %v, want value1", val)
	}
}

func TestMemStorageGetKeys(t *testing.T) {
	store := NewMemStorage[string]()
	store.Set("key1", "subkey1", "value1")
	store.Set("key1", "subkey2", "value2")
	store.Set("key2", "subkey3", "value3")

	keys, ok := store.GetKeys("key1")
	if !ok || len(keys) != 2 {
		t.Errorf("GetKeys() = %v, want 2 keys", keys)
	}
}

func TestMemStorageDelete(t *testing.T) {
	store := NewMemStorage[string]()
	store.Set("key1", "subkey1", "value1")
	store.Set("key1", "subkey2", "value2")

	store.Delete("key1", "subkey1")
	val, ok := store.Get("key1", "subkey1")
	if ok || val != "" {
		t.Errorf("Get() after Delete() = %v, want empty value", val)
	}
}

func TestMemStorageDeleteAll(t *testing.T) {
	store := NewMemStorage[string]()
	store.Set("key1", "subkey1", "value1")
	store.Set("key1", "subkey2", "value2")

	store.DeleteAll("key1")
	val, ok := store.Get("key1", "subkey1")
	if ok || val != "" {
		t.Errorf("Get() after DeleteAll() = %v, want empty value", val)
	}
}

func TestMemStorageExists(t *testing.T) {
	store := NewMemStorage[string]()
	store.Set("key1", "subkey1", "value1")

	if !store.Exists("key1") {
		t.Errorf("Exists() = false, want true")
	}

	if store.Exists("key2") {
		t.Errorf("Exists() = true, want false")
	}
}
