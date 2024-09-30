package net

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
	n := NewNet("testing")
	err := cache.Set(n)
	if err != nil {
		t.Error("failed to set net in cache")
	}
}

func TestGetCache(t *testing.T) {
	cache := NewCache()
	n := NewNet("testing")
	cache.Set(n)
	_, ok := cache.Get(n.ID)
	if !ok {
		t.Error("failed to get net from cache")
	}
}

func TestDeleteCache(t *testing.T) {
	cache := NewCache()
	n := NewNet("testing")
	cache.Set(n)
	cache.Delete(n.ID)
	_, ok := cache.Get(n.ID)
	if ok {
		t.Error("failed to delete net from cache")
	}
}

func TestCloseCache(t *testing.T) {
	cache := NewCache()
	n := NewNet("testing")
	cache.Set(n)
	cache.Close()
	_, ok := cache.Get(n.ID)
	if ok {
		t.Error("failed to close cache")
	}
}

func TestModelFileName(t *testing.T) {
	cache := NewCache()
	cache.ModelsDir = "models"
	model := "testing.onnx"
	modelFile := cache.ModelFileName(model)
	if modelFile != "models/testing.onnx" {
		t.Error("failed to get model file name")
	}
}

func TestKnownModelFileName(t *testing.T) {
	cache := NewCache()
	cache.ModelsDir = "models"
	model := "candy-9"
	modelFile := cache.ModelFileName(model)
	if modelFile != "models/candy-9.onnx" {
		t.Error("failed to get known model file name")
	}
}
