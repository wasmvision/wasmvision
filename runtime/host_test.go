package runtime

import (
	"context"
	"testing"
)

func TestNewInterpreter(t *testing.T) {
	interpreter := New(context.Background(), InterpreterConfig{})
	if len(interpreter.Processors()) != 0 {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_Processors(t *testing.T) {
	interpreter := New(context.Background(), InterpreterConfig{})
	if len(interpreter.Processors()) != 0 {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_LoadProcessors(t *testing.T) {
	interpreter := New(context.Background(), InterpreterConfig{})
	if err := interpreter.LoadProcessors(context.Background(), []string{}); err != nil {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_Close(t *testing.T) {
	interpreter := New(context.Background(), InterpreterConfig{})
	interpreter.Close(context.Background())
}
