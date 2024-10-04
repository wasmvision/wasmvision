package guest

import (
	"testing"
)

func TestWellKnownProcessor(t *testing.T) {
	t.Run("well-known processor", func(t *testing.T) {
		if !ProcessorWellKnown("candy") {
			t.Errorf("processor candy not found")
		}
	})

	t.Run("unknown processor", func(t *testing.T) {
		if ProcessorWellKnown("unknown") {
			t.Errorf("processor unknown found")
		}
	})
}
