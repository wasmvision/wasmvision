//go:build yzma

package runtime

import (
	"os"
	"testing"
)

func TestVLM(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("you must set YZMA_LIB to run this test")
	}

	if err := llamaInit(); err != nil {
		t.Fatal(err)
	}
	defer llamaFree()
}
