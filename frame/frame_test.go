package frame

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestNewFrame(t *testing.T) {
	frm := NewFrame()
	if frm.ID == 0 {
		t.Error("frame ID is 0")
	}
}

func TestSetImage(t *testing.T) {
	frm := NewFrame()
	img := gocv.NewMat()
	frm.SetImage(img)
	if !frm.Empty() {
		t.Error("frame image should be empty")
	}
}

func TestClose(t *testing.T) {
	frm := NewFrame()
	img := gocv.NewMat()
	frm.SetImage(img)
	frm.Close()
	if !frm.Image.Closed() {
		t.Error("frame image is not closed")
	}
}
