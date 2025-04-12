package cv

import (
	"testing"

	"github.com/orsinium-labs/tinytest/is"
	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

func TestBlur(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 480)
	is.Equal(c, frm.Raw.Image.Cols(), 640)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	sz := Size{X: wypes.Int32(10), Y: wypes.Int32(10)}

	res := testFrameResult()

	blurFunc(ctx)(&store, frm, sz, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestBoxFilter(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, 5)

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 480)
	is.Equal(c, frm.Raw.Image.Cols(), 640)
	is.Equal(c, frm.Raw.Image.Type(), gocv.MatTypeCV32F)

	sz := Size{X: wypes.Int32(10), Y: wypes.Int32(10)}

	res := testFrameResult()

	boxFilterFunc(ctx)(&store, frm, wypes.UInt32(6), sz, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestEstimateAffine2d(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm1 := matNewWithSizeFunc(ctx)(&store, 4, 1, wypes.UInt32(gocv.MatTypeCV32FC2))
	frm2 := matNewWithSizeFunc(ctx)(&store, 4, 1, wypes.UInt32(gocv.MatTypeCV32FC2))

	frm1.Raw.Image.SetFloatAt(0, 0, float32(0))
	frm1.Raw.Image.SetFloatAt(0, 1, float32(0))
	frm1.Raw.Image.SetFloatAt(1, 0, float32(10))
	frm1.Raw.Image.SetFloatAt(1, 1, float32(5))
	frm1.Raw.Image.SetFloatAt(2, 0, float32(10))
	frm1.Raw.Image.SetFloatAt(2, 1, float32(10))
	frm1.Raw.Image.SetFloatAt(3, 0, float32(5))
	frm1.Raw.Image.SetFloatAt(3, 1, float32(10))

	frm2.Raw.Image.SetFloatAt(0, 0, float32(0))
	frm2.Raw.Image.SetFloatAt(0, 1, float32(0))
	frm2.Raw.Image.SetFloatAt(1, 0, float32(10))
	frm2.Raw.Image.SetFloatAt(1, 1, float32(0))
	frm2.Raw.Image.SetFloatAt(2, 0, float32(10))
	frm2.Raw.Image.SetFloatAt(2, 1, float32(10))
	frm2.Raw.Image.SetFloatAt(3, 0, float32(0))
	frm2.Raw.Image.SetFloatAt(3, 1, float32(10))

	res := testFrameResult()

	estimateAffine2dFunc(ctx)(&store, frm1, frm2, res)

	result := getResult(&store, res)

	if result.IsError {
		t.Errorf("error: %v", result.Error)
		return
	}
	frm := result.OK.Raw
	if frm == nil {
		t.Errorf("error: nil frame")
		return
	}

	is.Equal(c, frm.Image.Rows(), 2)
	is.Equal(c, frm.Image.Cols(), 3)
}

func TestAdaptiveThreshold(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	is.Equal(c, frm.Raw.Empty(), false)
	is.Equal(c, frm.Raw.Image.Rows(), 480)
	is.Equal(c, frm.Raw.Image.Cols(), 640)

	res := testFrameResult()

	adaptiveThresholdFunc(ctx)(&store, frm, 255, wypes.UInt32(gocv.AdaptiveThresholdMean), wypes.UInt32(gocv.ThresholdBinary), 11, 2, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestAdd(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()
	ctx.ReturnDataPtr = 512

	frm1 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))
	frm2 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	res := testFrameResult()

	addFunc(ctx)(&store, frm1, frm2, res)

	store.Stack.Push(64)
	result := res.Lift(&store)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestSubtract(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm1 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))
	frm2 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	res := testFrameResult()

	subtractFunc(ctx)(&store, frm1, frm2, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestMultiply(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm1 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))
	frm2 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	res := testFrameResult()

	multiplyFunc(ctx)(&store, frm1, frm2, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestDivide(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm1 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))
	frm2 := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	res := testFrameResult()

	divideFunc(ctx)(&store, frm1, frm2, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestExp(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV32FC1))

	res := testFrameResult()

	expFunc(ctx)(&store, frm, res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestGaussianBlur(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	sz := Size{X: wypes.Int32(5), Y: wypes.Int32(5)}

	res := testFrameResult()

	gaussianBlurFunc(ctx)(&store, frm, sz, 1.5, 1.5, wypes.UInt32(gocv.BorderDefault), res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}

func TestNormalize(t *testing.T) {
	c := is.NewRelaxed(t)
	store := getTestStore()
	ctx := getTestCVContext()

	frm := matNewWithSizeFunc(ctx)(&store, 480, 640, wypes.UInt32(gocv.MatTypeCV8UC1))

	res := testFrameResult()

	normalizeFunc(ctx)(&store, frm, wypes.Float32(2), wypes.Float32(2), wypes.UInt32(2), res)

	result := getResult(&store, res)

	is.Equal(c, result.OK.Raw.Image.Cols(), 640)
	is.Equal(c, result.OK.Raw.Image.Rows(), 480)
}
