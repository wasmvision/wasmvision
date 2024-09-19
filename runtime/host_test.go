package runtime

import (
	"context"
	"testing"
)

func TestNewInterpreter(t *testing.T) {
	interpreter := New(context.Background())
	if len(interpreter.Processors()) != 0 {
		t.Error("interpreter is invalid")
	}
}
