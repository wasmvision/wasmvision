package cv

import (
	"testing"

	"github.com/orsinium-labs/tinytest/is"
	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func TestNewMat(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewFunc(ctx)(&store)

	is.Equal(c, frm.Raw.Empty(), true)
}

func TestNewMatWithSize(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 16)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV8UC3)
}

func TestAddFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	matAddFloatFunc(ctx)(&store, frm, wypes.Float32(1.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 1.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 1.0)
}

func TestSubtractFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	matSubtractFloatFunc(ctx)(&store, frm, wypes.Float32(1.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), -1.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), -1.0)
}

func TestMultiplyFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	frm.Raw.Image.SetFloatAt(1, 1, 1.0)
	frm.Raw.Image.SetFloatAt(320, 240, 1.0)

	matMultiplyFloatFunc(ctx)(&store, frm, wypes.Float32(2.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 2.0)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 2.0)
}

func TestDivideFloat(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 640)
	is.Equal(c, frm.Raw.Image.Cols(), 480)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	frm.Raw.Image.SetFloatAt(1, 1, 1.0)
	frm.Raw.Image.SetFloatAt(320, 240, 1.0)

	matDivideFloatFunc(ctx)(&store, frm, wypes.Float32(2.0))

	is.Equal(c, frm.Raw.Image.GetFloatAt(1, 1), 0.5)
	is.Equal(c, frm.Raw.Image.GetFloatAt(320, 240), 0.5)
}

func TestMatClone(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	cloned := matCloneFunc(ctx)(&store, frm)

	is.Equal(c, cloned.Raw.Empty(), false)
	is.Equal(c, cloned.Raw.Image.Rows(), 640)
	is.Equal(c, cloned.Raw.Image.Cols(), 480)
	is.Equal(c, cloned.Raw.Image.Type(), gocv.MatTypeCV32F)
}

func TestMatColsRows(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	cols := matColsFunc(ctx)(&store, frm)
	is.Equal(c, cols, wypes.UInt32(480))

	rows := matRowsFunc(ctx)(&store, frm)
	is.Equal(c, rows, wypes.UInt32(640))
}

func TestMatRegion(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	rect := Rect{
		Min: Size{X: wypes.Int32(100), Y: wypes.Int32(100)},
		Max: Size{X: wypes.Int32(300), Y: wypes.Int32(300)},
	}

	region := matRegionFunc(ctx)(&store, frm, rect)

	is.Equal(c, region.Raw.Image.Rows(), 200)
	is.Equal(c, region.Raw.Image.Cols(), 200)
}

func TestMatEmpty(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)
	is.Equal(c, matEmptyFunc(ctx)(&store, frm), wypes.Bool(false))
}

func TestMatType(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	matType := matTypeFunc(ctx)(&store, frm)

	is.Equal(c, matType, wypes.UInt32(gocv.MatTypeCV32F))
}

func TestMatSize(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	sizeList := wypes.ReturnedList[wypes.UInt32]{
		DataPtr: 512,
		Offset:  64,
	}

	matSizeFunc(ctx)(&store, frm, sizeList)

	store.Stack.Push(64)
	result := sizeList.Lift(&store)

	is.Equal(c, len(result.Raw), 2)
	is.Equal(c, result.Raw[0], wypes.UInt32(640))
	is.Equal(c, result.Raw[1], wypes.UInt32(480))
}

func TestMatReshape(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	res := testFrameResult()

	matReshapeFunc(ctx)(&store, frm, 1, 307200, res)

	store.Stack.Push(64)
	result := res.Lift(&store)

	is.Equal(c, result.OK.Raw.Image.Rows(), 307200)
	is.Equal(c, result.OK.Raw.Image.Cols(), 1)
}

func TestMatGetSetFloatAt(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	matSetFloatAtFunc(ctx)(&store, frm, 10, 10, wypes.Float32(3.14))
	value := matGetFloatAtFunc(ctx)(&store, frm, 10, 10)

	is.Equal(c, value, wypes.Float32(3.14))
}

func TestMatGetSetUcharAt(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	matSetUcharAtFunc(ctx)(&store, frm, 20, 20, wypes.UInt8(255))
	value := matGetUcharAtFunc(ctx)(&store, frm, 20, 20)

	is.Equal(c, value, wypes.UInt8(255))
}

func TestMatGetVecbAt(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 16)

	is.Equal(c, frm.Raw.Empty(), false)

	vecbList := wypes.ReturnedList[wypes.UInt32]{
		DataPtr: 512,
		Offset:  64,
	}

	matGetVecbAtFunc(ctx)(&store, frm, 30, 30, vecbList)

	store.Stack.Push(64)
	result := vecbList.Lift(&store)

	is.Equal(c, len(result.Raw), 3)
	is.Equal(c, result.Raw[0], wypes.UInt32(0))
	is.Equal(c, result.Raw[1], wypes.UInt32(0))
	is.Equal(c, result.Raw[2], wypes.UInt32(0))
}

func TestMatMinMaxLoc(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 640, 480, 5)

	is.Equal(c, frm.Raw.Empty(), false)

	frm.Raw.Image.SetFloatAt(0, 0, 1.0)
	frm.Raw.Image.SetFloatAt(639, 479, 100.0)

	res := wypes.Result[MixMaxLocResult, MixMaxLocResult, wypes.UInt32]{
		DataPtr: 512,
		Offset:  64,
	}

	matMinMaxLocFunc(ctx)(&store, frm, res)

	store.Stack.Push(64)
	result := res.Lift(&store)

	is.Equal(c, result.OK.MinVal, wypes.Float32(0))
	is.Equal(c, result.OK.MaxVal, wypes.Float32(100.0))
	is.Equal(c, result.OK.MinLoc.X, wypes.Int32(1.0))
	is.Equal(c, result.OK.MinLoc.Y, wypes.Int32(0))
	is.Equal(c, result.OK.MaxLoc.X, wypes.Int32(479))
	is.Equal(c, result.OK.MaxLoc.Y, wypes.Int32(639))
}

func TestMatOnes(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	res := testFrameResult()

	matOnesFunc(ctx)(&store, 640, 480, 5, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Empty(), false)
	is.Equal(c, result.OK.Raw.Image.Rows(), 640)
	is.Equal(c, result.OK.Raw.Image.Cols(), 480)
	is.Equal(c, result.OK.Raw.Image.GetFloatAt(0, 0), 1.0)
}

func TestMatZeros(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	res := testFrameResult()

	matZerosFunc(ctx)(&store, 640, 480, 5, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Empty(), false)
	is.Equal(c, result.OK.Raw.Image.Rows(), 640)
	is.Equal(c, result.OK.Raw.Image.Cols(), 480)
	is.Equal(c, result.OK.Raw.Image.GetFloatAt(0, 0), 0.0)
}
