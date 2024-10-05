package net

import (
	"testing"
)

func TestNewNet(t *testing.T) {
	t.Run("new net", func(t *testing.T) {
		n := NewNet("model")
		if n.Name != "model" {
			t.Errorf("unexpected model name %s", n.Name)
		}
	})
}
