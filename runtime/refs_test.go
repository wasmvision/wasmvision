package runtime

import (
	"testing"

	"github.com/orsinium-labs/tinytest/is"
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
)

func TestHostRef_Lower(t *testing.T) {
	c := is.NewRelaxed(t)
	stack := wypes.NewSliceStack(4)
	store := wypes.Store{
		Stack: stack,
		Refs:  wypes.NewMapRefs(),
	}

	frm1 := cv.NewEmptyFrame()
	frm1.ID = 1
	frm2 := cv.NewEmptyFrame()
	frm2.ID = 2

	defer frm1.Close()
	defer frm2.Close()

	val1 := wypes.HostRef[*cv.Frame]{Raw: frm1}
	val2 := wypes.HostRef[*cv.Frame]{Raw: frm2}
	val1.Lower(&store)
	val2.Lower(&store)
	val3 := val2.Lift(&store)
	is.Equal(c, val3.Unwrap().ID, frm2.ID)
}

func TestHostRef_Drop(t *testing.T) {
	c := is.NewRelaxed(t)
	refs := wypes.NewMapRefs()
	store := wypes.Store{
		Stack: wypes.NewSliceStack(4),
		Refs:  refs,
	}
	frm1 := cv.NewEmptyFrame()
	frm1.ID = 1

	defer frm1.Close()

	val1 := wypes.HostRef[*cv.Frame]{Raw: frm1}
	val1.Lower(&store)
	val2 := val1.Lift(&store)
	is.Equal(c, val2.Unwrap().ID, frm1.ID)
	is.Equal(c, len(refs.Raw), 1)
	val2.Drop()
	is.Equal(c, len(refs.Raw), 0)
}
