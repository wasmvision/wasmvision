package cv

import (
	"image"
	"image/color"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
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
	var T wypes.Int32
	return Size{
		Y: T.Lift(s),
		X: T.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v Size) Lower(s *wypes.Store) {
	v.Y.Lower(s)
	v.X.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Size) MemoryLift(s *wypes.Store, offset uint32) (Size, uint32) {
	var T wypes.Int32

	y, ySize := T.MemoryLift(s, offset)
	x, xSize := T.MemoryLift(s, offset+ySize)

	return Size{X: x, Y: y}, xSize + ySize
}

// MemoryLower implements [MemoryLower] interface.
func (v Size) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	ySize := v.Y.MemoryLower(s, offset)
	xSize := v.X.MemoryLower(s, offset+ySize)

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
	var T Size
	return Rect{
		Min: T.Lift(s),
		Max: T.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v Rect) Lower(s *wypes.Store) {
	v.Min.Lower(s)
	v.Max.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Rect) MemoryLift(s *wypes.Store, offset uint32) (Rect, uint32) {
	var T Size

	min, minSize := T.MemoryLift(s, offset)
	max, maxSize := T.MemoryLift(s, offset+minSize)

	return Rect{Min: min, Max: max}, minSize + maxSize
}

// MemoryLower implements [MemoryLower] interface.
func (v Rect) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	minSize := v.Min.MemoryLower(s, offset)
	maxSize := v.Max.MemoryLower(s, offset+minSize)

	return minSize + maxSize
}

// Scalar represents the record "wasm:cv/types#scalar".
//
//	record scalar {
//		val1: f32,
//		val2: f32,
//		val3: f32,
//		val4: f32,
//	}
type Scalar struct {
	Val1 wypes.Float32
	Val2 wypes.Float32
	Val3 wypes.Float32
	Val4 wypes.Float32
}

func (v Scalar) Unwrap() gocv.Scalar {
	return gocv.NewScalar(float64(v.Val1.Unwrap()), float64(v.Val2.Unwrap()), float64(v.Val3.Unwrap()), float64(v.Val4.Unwrap()))
}

func (v Scalar) ValueTypes() []wypes.ValueType {
	return []wypes.ValueType{wypes.ValueTypeF32, wypes.ValueTypeF32, wypes.ValueTypeF32, wypes.ValueTypeF32}
}

func (Scalar) Lift(s *wypes.Store) Scalar {
	var val1 wypes.Float32
	var val2 wypes.Float32
	var val3 wypes.Float32
	var val4 wypes.Float32
	return Scalar{
		Val1: val1.Lift(s),
		Val2: val2.Lift(s),
		Val3: val3.Lift(s),
		Val4: val4.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v Scalar) Lower(s *wypes.Store) {
	v.Val1.Lower(s)
	v.Val2.Lower(s)
	v.Val3.Lower(s)
	v.Val4.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Scalar) MemoryLift(s *wypes.Store, offset uint32) (Scalar, uint32) {
	var T wypes.Float32

	val1, v1Size := T.MemoryLift(s, offset)
	val2, v2Size := T.MemoryLift(s, offset+v1Size)
	val3, v3Size := T.MemoryLift(s, offset+v1Size+v2Size)
	val4, v4Size := T.MemoryLift(s, offset+v1Size+v2Size+v3Size)

	return Scalar{
		Val1: val1,
		Val2: val2,
		Val3: val3,
		Val4: val4,
	}, v1Size + v2Size + v3Size + v4Size
}

// MemoryLower implements [MemoryLower] interface.
func (v Scalar) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	v1Size := v.Val1.MemoryLower(s, offset)
	v2Size := v.Val2.MemoryLower(s, offset+v1Size)
	v3Size := v.Val3.MemoryLower(s, offset+v1Size+v2Size)
	v4Size := v.Val4.MemoryLower(s, offset+v1Size+v2Size+v3Size)

	return v1Size + v2Size + v3Size + v4Size
}

// RGBA represents the record "wasm:cv/types#RGBA".
//
//	record RGBA {
//		r: u8,
//		g: u8,
//		b: u8,
//		a: u8,
//	}
type RGBA struct {
	R wypes.UInt8
	G wypes.UInt8
	B wypes.UInt8
	A wypes.UInt8
}

func (v RGBA) Unwrap() color.RGBA {
	return color.RGBA{R: v.R.Unwrap(), G: v.G.Unwrap(), B: v.B.Unwrap(), A: v.A.Unwrap()}
}

func (v RGBA) ValueTypes() []wypes.ValueType {
	return []wypes.ValueType{wypes.ValueTypeI32, wypes.ValueTypeI32, wypes.ValueTypeI32, wypes.ValueTypeI32}
}

func (RGBA) Lift(s *wypes.Store) RGBA {
	var r wypes.UInt8
	var g wypes.UInt8
	var b wypes.UInt8
	var a wypes.UInt8
	return RGBA{
		R: r.Lift(s),
		G: g.Lift(s),
		B: b.Lift(s),
		A: a.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v RGBA) Lower(s *wypes.Store) {
	v.R.Lower(s)
	v.G.Lower(s)
	v.B.Lower(s)
	v.A.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (RGBA) MemoryLift(s *wypes.Store, offset uint32) (RGBA, uint32) {
	var T wypes.UInt8

	r, rSize := T.MemoryLift(s, offset)
	g, gSize := T.MemoryLift(s, offset+rSize)
	b, bSize := T.MemoryLift(s, offset+rSize+gSize)
	a, aSize := T.MemoryLift(s, offset+rSize+gSize+bSize)

	return RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}, rSize + gSize + bSize + aSize
}

// MemoryLower implements [MemoryLower] interface.
func (v RGBA) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	rSize := v.R.MemoryLower(s, offset)
	gSize := v.G.MemoryLower(s, offset+rSize)
	bSize := v.B.MemoryLower(s, offset+rSize+gSize)
	aSize := v.A.MemoryLower(s, offset+rSize+gSize+bSize)

	return rSize + gSize + bSize + aSize
}
