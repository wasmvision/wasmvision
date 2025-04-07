package cv

import (
	"testing"

	"github.com/orsinium-labs/tinytest/is"
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/config"
	"gocv.io/x/gocv"
)

func TestNewMat(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewFunc(ctx)
	frm := f(&store)

	is.Equal(c, frm.Raw.Empty(), true)
}

func TestNewMatWithSize(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewWithSizeFunc(ctx)
	frm := f(&store, 640, 480, 16)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV8UC3)
}

func TestAddFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewWithSizeFunc(ctx)
	frm := f(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	f2 := matAddFloatFunc(ctx)
	f2(&store, frm, wypes.Float32(1.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 1.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 1.0)
}

func TestSubtractFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewWithSizeFunc(ctx)
	frm := f(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	f2 := matSubtractFloatFunc(ctx)
	f2(&store, frm, wypes.Float32(1.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), -1.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), -1.0)
}

func TestMultiplyFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewWithSizeFunc(ctx)
	frm := f(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	frm.Raw.Image.SetFloatAt(1, 1, 1.0)
	frm.Raw.Image.SetFloatAt(320, 240, 1.0)

	f2 := matMultiplyFloatFunc(ctx)
	f2(&store, frm, wypes.Float32(2.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 2.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 2.0)
}

func TestDivideFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(1024),
		Refs:   wypes.NewMapRefs(),
	}

	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("", configStore, false)

	f := matNewWithSizeFunc(ctx)
	frm := f(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	frm.Raw.Image.SetFloatAt(1, 1, 1.0)
	frm.Raw.Image.SetFloatAt(320, 240, 1.0)

	f2 := matDivideFloatFunc(ctx)
	f2(&store, frm, wypes.Float32(2.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 0.5)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 0.5)
}
