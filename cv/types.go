package cv

import (
	"image"

	"github.com/orsinium-labs/wypes"
)

// Size represents the record "wasm:cv/types#size".
//
//	record size {
//		x: s32,
//		y: s32,
//	}
type Size struct {
	X wypes.Int32
	Y wypes.Int32
}

func (v Size) Unwrap() image.Point {
	return image.Point{X: int(v.X.Unwrap()), Y: int(v.Y.Unwrap())}
}

func (v Size) ValueTypes() []wypes.ValueType {
	types := make([]wypes.ValueType, 0, 2)
	types = append(types, v.X.ValueTypes()...)
	types = append(types, v.Y.ValueTypes()...)
	return types
}

func (Size) Lift(s *wypes.Store) Size {
	var x wypes.Int32
	var y wypes.Int32
	return Size{
		X: x.Lift(s),
		Y: y.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v Size) Lower(s *wypes.Store) {
	v.X.Lower(s)
	v.Y.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Size) MemoryLift(s *wypes.Store, offset uint32) (Size, uint32) {
	var x wypes.Int32
	var y wypes.Int32

	x, xSize := x.MemoryLift(s, offset)
	y, ySize := y.MemoryLift(s, offset+xSize)

	return Size{X: x, Y: y}, xSize + ySize
}

// MemoryLower implements [MemoryLower] interface.
func (v Size) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	xSize := v.X.MemoryLower(s, offset)
	ySize := v.Y.MemoryLower(s, offset+xSize)

	return xSize + ySize
}

// Rect represents the record "wasm:cv/types#rect".
//
//	record rect {
//		min: size,
//		max: size,
//	}
type Rect struct {
	Min Size
	Max Size
}

func (v Rect) Unwrap() image.Rectangle {
	return image.Rectangle{Min: image.Point{X: int(v.Min.X.Unwrap()), Y: int(v.Min.Y.Unwrap())}, Max: image.Point{X: int(v.Max.X.Unwrap()), Y: int(v.Max.Y.Unwrap())}}
}

func (v Rect) ValueTypes() []wypes.ValueType {
	types := make([]wypes.ValueType, 0, 4)
	types = append(types, v.Min.ValueTypes()...)
	types = append(types, v.Max.ValueTypes()...)
	return types
}

func (Rect) Lift(s *wypes.Store) Rect {
	var min Size
	var max Size
	return Rect{
		Min: min.Lift(s),
		Max: max.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v Rect) Lower(s *wypes.Store) {
	v.Min.Lower(s)
	v.Max.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Rect) MemoryLift(s *wypes.Store, offset uint32) (Rect, uint32) {
	var min Size
	var max Size

	min, minSize := min.MemoryLift(s, offset)
	max, maxSize := max.MemoryLift(s, offset+minSize)

	return Rect{Min: min, Max: max}, minSize + maxSize
}

// MemoryLower implements [MemoryLower] interface.
func (v Rect) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	minSize := v.Min.MemoryLower(s, offset)
	maxSize := v.Max.MemoryLower(s, offset+minSize)

	return minSize + maxSize
}
