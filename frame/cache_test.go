package frame

import (
	"testing"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()
	if cache == nil {
		t.Error("cache is nil")
	}
}

func TestSetCache(t *testing.T) {
	cache := NewCache()
	frm := NewFrame()
	err := cache.Set(frm)
	if err != nil {
		t.Error("failed to set frame in cache")
	}
}

func TestGetCache(t *testing.T) {
	cache := NewCache()
	frm := NewFrame()
	cache.Set(frm)
	_, ok := cache.Get(frm.ID)
	if !ok {
		t.Error("failed to get frame from cache")
	}
}

func TestDeleteCache(t *testing.T) {
	cache := NewCache()
	frm := NewFrame()
	cache.Set(frm)
	cache.Delete(frm.ID)
	_, ok := cache.Get(frm.ID)
	if ok {
		t.Error("failed to delete frame from cache")
	}
}

func TestCloseCache(t *testing.T) {
	cache := NewCache()
	frm := NewFrame()
	cache.Set(frm)
	cache.Close()
	_, ok := cache.Get(frm.ID)
	if ok {
		t.Error("failed to close cache")
	}
}
