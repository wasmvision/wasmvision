package runtime

import (
	"context"
	"os"
	"testing"
)

func TestNewInterpreter(t *testing.T) {
	interpreter, _ := New(context.Background(), InterpreterConfig{})
	if len(interpreter.Processors()) != 0 {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_Processors(t *testing.T) {
	interpreter, _ := New(context.Background(), InterpreterConfig{})
	if len(interpreter.Processors()) != 0 {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_LoadProcessors(t *testing.T) {
	interpreter, _ := New(context.Background(), InterpreterConfig{})
	if err := interpreter.LoadProcessors(context.Background(), []string{}); err != nil {
		t.Error("interpreter is invalid")
	}
}

func TestInterpreter_Close(t *testing.T) {
	interpreter, _ := New(context.Background(), InterpreterConfig{})
	interpreter.Close(context.Background())
}

func TestInterpreter_RegisterGuestModule(t *testing.T) {
	module, err := os.ReadFile("../processors/hello.wasm")
	if err != nil {
		t.Error("failed to read hello.wasm")
	}

	interpreter, _ := New(context.Background(), InterpreterConfig{})
	if err := interpreter.RegisterGuestModule(context.Background(), "hello", module); err != nil {
		t.Error("interpreter is invalid")
	}
}

func TestLoadProcessors(t *testing.T) {
	interpreter, _ := New(context.Background(), InterpreterConfig{})
	if err := interpreter.LoadProcessors(context.Background(), []string{"../processors/hello.wasm"}); err != nil {
		t.Error("interpreter is invalid")
	}
}
