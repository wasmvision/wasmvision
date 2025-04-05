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
	v.X.Lower(s)
	v.Y.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (Size) MemoryLift(s *wypes.Store, offset uint32) (Size, uint32) {
	var T wypes.Int32

	x, xSize := T.MemoryLift(s, offset)
	y, ySize := T.MemoryLift(s, offset+xSize)

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
	types := make([]wypes.ValueType, 0)
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
	v.Max.Lower(s)
	v.Min.Lower(s)
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
		Val4: val4.Lift(s),
		Val3: val3.Lift(s),
		Val2: val2.Lift(s),
		Val1: val1.Lift(s),
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
	offset += v1Size

	val2, v2Size := T.MemoryLift(s, offset)
	offset += v2Size

	val3, v3Size := T.MemoryLift(s, offset)
	offset += v3Size

	val4, v4Size := T.MemoryLift(s, offset)
	offset += v4Size

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
	offset += v1Size

	v2Size := v.Val2.MemoryLower(s, offset)
	offset += v2Size

	v3Size := v.Val3.MemoryLower(s, offset)
	offset += v3Size

	v4Size := v.Val4.MemoryLower(s, offset)
	offset += v4Size

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
		A: a.Lift(s),
		B: b.Lift(s),
		G: g.Lift(s),
		R: r.Lift(s),
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

// MixMaxLocResult represents the record "wasm:cv/types#mix-max-loc-result".
//
//	record mix-max-loc-result {
//		min-val: f32,
//		max-val: f32,
//		min-loc: size,
//		max-loc: size,
//	}
type MixMaxLocResult struct {
	MinVal wypes.Float32
	MaxVal wypes.Float32
	MinLoc Size
	MaxLoc Size
}

func (v MixMaxLocResult) Unwrap() MixMaxLocResult {
	return v
}

func (v MixMaxLocResult) ValueTypes() []wypes.ValueType {
	types := make([]wypes.ValueType, 0)
	types = append(types, v.MinVal.ValueTypes()...)
	types = append(types, v.MaxVal.ValueTypes()...)
	types = append(types, v.MinLoc.ValueTypes()...)
	types = append(types, v.MaxLoc.ValueTypes()...)
	return types
}

func (MixMaxLocResult) Lift(s *wypes.Store) MixMaxLocResult {
	var (
		T Size
		V wypes.Float32
	)
	return MixMaxLocResult{
		MinVal: V.Lift(s),
		MaxVal: V.Lift(s),
		MinLoc: T.Lift(s),
		MaxLoc: T.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v MixMaxLocResult) Lower(s *wypes.Store) {
	v.MinVal.Lower(s)
	v.MaxVal.Lower(s)
	v.MinLoc.Lower(s)
	v.MaxLoc.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (MixMaxLocResult) MemoryLift(s *wypes.Store, offset uint32) (MixMaxLocResult, uint32) {
	var (
		T Size
		V wypes.Float32
	)

	minVal, minValSize := V.MemoryLift(s, offset)
	offset += minValSize

	maxVal, maxValSize := V.MemoryLift(s, offset)
	offset += maxValSize

	minLoc, minSize := T.MemoryLift(s, offset)
	offset += minSize

	maxLoc, maxSize := T.MemoryLift(s, offset)
	offset += maxSize

	return MixMaxLocResult{MinVal: minVal, MaxVal: maxVal, MinLoc: minLoc, MaxLoc: maxLoc}, minSize + maxSize + maxValSize + minValSize
}

// MemoryLower implements [MemoryLower] interface.
func (v MixMaxLocResult) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	minValSize := v.MinVal.MemoryLower(s, offset)
	offset += minValSize

	maxValSize := v.MaxVal.MemoryLower(s, offset)
	offset += maxValSize

	minSize := v.MinLoc.MemoryLower(s, offset)
	offset += minSize

	maxSize := v.MaxLoc.MemoryLower(s, offset)
	offset += maxSize

	return minSize + maxSize + maxValSize + minValSize
}

// BlobParams represents the record "wasm:cv/types#blob-params".
//
//	record blob-params {
//		scale-factor: f32,
//		size: size,
//		mean: scalar,
//		swap-RB: bool,
//		ddepth: u8,
//		data-layout: data-layout-type,
//		padding-mode: padding-mode-type,
//		border: scalar,
//	}
type BlobParams struct {
	ScaleFactor wypes.Float32
	Size        Size
	Mean        Scalar
	SwapRB      wypes.Bool
	DDepth      wypes.UInt8
	DataLayout  wypes.UInt8
	PaddingMode wypes.UInt8
	Border      Scalar
}

func (v BlobParams) Unwrap() gocv.ImageToBlobParams {
	return gocv.ImageToBlobParams{
		ScaleFactor: float64(v.ScaleFactor.Unwrap()),
		Size:        v.Size.Unwrap(),
		Mean:        v.Mean.Unwrap(),
		SwapRB:      v.SwapRB.Unwrap(),
		Ddepth:      gocv.MatType(v.DDepth.Unwrap()),
		DataLayout:  gocv.DataLayoutType(v.DataLayout.Unwrap()),
		PaddingMode: gocv.PaddingModeType(v.PaddingMode.Unwrap()),
		BorderValue: v.Border.Unwrap(),
	}
}

func (v BlobParams) ValueTypes() []wypes.ValueType {
	types := make([]wypes.ValueType, 0)
	types = append(types, wypes.ValueTypeF32)
	types = append(types, v.Size.ValueTypes()...)
	types = append(types, v.Mean.ValueTypes()...)
	types = append(types, v.SwapRB.ValueTypes()...)
	types = append(types, v.DDepth.ValueTypes()...)
	types = append(types, v.DataLayout.ValueTypes()...)
	types = append(types, v.PaddingMode.ValueTypes()...)
	types = append(types, v.Border.ValueTypes()...)
	return types
}

func (BlobParams) Lift(s *wypes.Store) BlobParams {
	var (
		T Size
		S Scalar
		V wypes.Float32
		B wypes.Bool
		U wypes.UInt8
	)
	return BlobParams{
		Border:      S.Lift(s),
		PaddingMode: U.Lift(s),
		DataLayout:  U.Lift(s),
		DDepth:      U.Lift(s),
		SwapRB:      B.Lift(s),
		Mean:        S.Lift(s),
		Size:        T.Lift(s),
		ScaleFactor: V.Lift(s),
	}
}

// Lower implements [Lower] interface.
func (v BlobParams) Lower(s *wypes.Store) {
	v.Border.Lower(s)
	v.PaddingMode.Lower(s)
	v.DataLayout.Lower(s)
	v.DDepth.Lower(s)
	v.SwapRB.Lower(s)
	v.Mean.Lower(s)
	v.Size.Lower(s)
	v.ScaleFactor.Lower(s)
}

// MemoryLift implements [MemoryLift] interface.
func (BlobParams) MemoryLift(s *wypes.Store, offset uint32) (BlobParams, uint32) {
	var (
		T Size
		S Scalar
		V wypes.Float32
		B wypes.Bool
		U wypes.UInt8
	)

	start := offset
	scaleFactor, scaleFactorSize := V.MemoryLift(s, offset)
	offset += scaleFactorSize
	size, sizeSize := T.MemoryLift(s, offset)
	offset += sizeSize
	mean, meanSize := S.MemoryLift(s, offset)
	offset += meanSize
	swapRB, swapRBSize := B.MemoryLift(s, offset)
	offset += swapRBSize
	ddepth, ddepthSize := U.MemoryLift(s, offset)
	offset += ddepthSize
	dataLayout, dataLayoutSize := U.MemoryLift(s, offset)
	offset += dataLayoutSize
	paddingMode, paddingModeSize := U.MemoryLift(s, offset)
	offset += paddingModeSize
	border, borderSize := S.MemoryLift(s, offset)
	offset += borderSize

	return BlobParams{
		ScaleFactor: scaleFactor,
		Size:        size,
		Mean:        mean,
		SwapRB:      swapRB,
		DDepth:      ddepth,
		DataLayout:  dataLayout,
		PaddingMode: paddingMode,
		Border:      border,
	}, offset - start
}

// MemoryLower implements [MemoryLower] interface.
func (v BlobParams) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	start := offset
	scaleFactorSize := v.ScaleFactor.MemoryLower(s, offset)
	offset += scaleFactorSize
	sizeSize := v.Size.MemoryLower(s, offset)
	offset += sizeSize
	meanSize := v.Mean.MemoryLower(s, offset)
	offset += meanSize
	swapRBSize := v.SwapRB.MemoryLower(s, offset)
	offset += swapRBSize
	ddepthSize := v.DDepth.MemoryLower(s, offset)
	offset += ddepthSize
	dataLayoutSize := v.DataLayout.MemoryLower(s, offset)
	offset += dataLayoutSize
	paddingModeSize := v.PaddingMode.MemoryLower(s, offset)
	offset += paddingModeSize
	borderSize := v.Border.MemoryLower(s, offset)
	offset += borderSize

	return offset - start
}

type BlobRectImageParams struct {
	Offset uint32
	Params BlobParams
	Rects  wypes.List[Rect]
	Size   Size
}

func (v BlobRectImageParams) Unwrap() BlobRectImageParams {
	return BlobRectImageParams{
		Params: v.Params,
		Rects:  v.Rects,
		Size:   v.Size,
	}
}

func (v BlobRectImageParams) ValueTypes() []wypes.ValueType {
	return []wypes.ValueType{wypes.ValueTypeI32}
}

func (BlobRectImageParams) Lift(s *wypes.Store) BlobRectImageParams {
	start := uint32(s.Stack.Pop())

	var (
		B BlobParams
		R wypes.List[Rect]
		S Size
	)
	offset := start

	params, paramsSize := B.MemoryLift(s, offset)
	offset += paramsSize
	list, _ := R.MemoryLift(s, offset)
	offset += 8 // size of the list pointer + list size
	size, sizeSize := S.MemoryLift(s, offset)
	offset += sizeSize

	return BlobRectImageParams{
		Offset: start,
		Params: params,
		Rects:  list,
		Size:   size,
	}
}

// Lower implements [Lower] interface.
func (v BlobRectImageParams) Lower(s *wypes.Store) {
	offset := v.Offset

	paramsSize := v.Params.MemoryLower(s, offset)
	offset += paramsSize
	listSize := v.Rects.MemoryLower(s, offset)
	offset += listSize
	sizeSize := v.Size.MemoryLower(s, offset)
	offset += sizeSize

	s.Stack.Push(wypes.Raw(v.Offset))
}

// MemoryLift implements [MemoryLift] interface.
func (BlobRectImageParams) MemoryLift(s *wypes.Store, offset uint32) (BlobRectImageParams, uint32) {
	var (
		B BlobParams
		R wypes.List[Rect]
		S Size
	)

	start := offset

	size, sizeSize := S.MemoryLift(s, offset)
	offset += sizeSize
	list, listSize := R.MemoryLift(s, offset)
	offset += listSize
	params, paramsSize := B.MemoryLift(s, offset)
	offset += paramsSize

	return BlobRectImageParams{
		Offset: start,
		Params: params,
		Rects:  list,
		Size:   size,
	}, paramsSize + listSize + sizeSize
}

// MemoryLower implements [MemoryLower] interface.
func (v BlobRectImageParams) MemoryLower(s *wypes.Store, offset uint32) (length uint32) {
	sizeSize := v.Size.MemoryLower(s, offset)
	offset += sizeSize
	listSize := v.Rects.MemoryLower(s, offset)
	offset += listSize
	paramsSize := v.Params.MemoryLower(s, offset)
	offset += paramsSize

	return paramsSize + listSize + sizeSize
}
