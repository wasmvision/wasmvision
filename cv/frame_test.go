package cv

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestNewFrame(t *testing.T) {
	frm := NewFrame(gocv.NewMat())
	defer frm.Close()

	if !frm.Image.Empty() {
		t.Error("frame image should be empty")
	}
}

func TestClose(t *testing.T) {
	frm := NewFrame(gocv.NewMat())
	if frm.Image.Closed() {
		t.Error("frame image is closed")
	}

	frm.Close()
	if !frm.Image.Closed() {
		t.Error("frame image is not closed")
	}
}
