package cv

import (
	"fmt"
	"log/slog"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"
)

func ObjDetectModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/objdetect": wypes.Module{
			"[constructor]cascade-classifier":               wypes.H2(newCascadeClassifierFunc(ctx)),
			"[resource-drop]cascade-classifier":             wypes.H2(closeFaceDetectorYNFunc(ctx)),
			"[method]cascade-classifier.close":              wypes.H2(closeCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.load":               wypes.H3(loadCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.detect-multi-scale": wypes.H4(detectMultiScaleCascadeClassifierFunc(ctx)),
			"[constructor]face-detector-YN":                 wypes.H4(newFaceDetectorYNFunc(ctx)),
			"[resource-drop]face-detector-YN":               wypes.H2(closeFaceDetectorYNFunc(ctx)),
			"[method]face-detector-YN.close":                wypes.H2(closeFaceDetectorYNFunc(ctx)),
			"[method]face-detector-YN.set-input-size":       wypes.H3(faceDetectorYNSetInputSizeFunc(ctx)),
			"[method]face-detector-YN.detect":               wypes.H3(faceDetectorYNDetectFunc(ctx)),
		},
	}
}

func newCascadeClassifierFunc[T *CascadeClassifier](ctx *Context) func(*wypes.Store, wypes.String) wypes.HostRef[T] {
	return func(s *wypes.Store, name wypes.String) wypes.HostRef[T] {
		c := gocv.NewCascadeClassifier()

		cl := NewCascadeClassifier(name.Raw)
		cl.SetClassifier(c)

		v := wypes.HostRef[T]{Raw: cl}
		return v
	}
}

func closeCascadeClassifierFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*CascadeClassifier]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*CascadeClassifier]) wypes.Void {
		cl := ref.Raw
		cl.Close()

		return wypes.Void{}
	}
}

func loadCascadeClassifierFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*CascadeClassifier], wypes.String) wypes.Bool {
	return func(s *wypes.Store, ref wypes.HostRef[*CascadeClassifier], file wypes.String) wypes.Bool {
		cl := ref.Raw
		name := file.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			slog.Info(fmt.Sprintf("Downloading classifier %s...", name))

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				slog.Error(fmt.Sprintf("Error downloading classifier: %v", err))
				return wypes.Bool(false)
			}

		case !models.ModelExists(modelFile):
			return wypes.Bool(false)
		}

		if cl == nil {
			slog.Error("classifier is nil")
			return wypes.Bool(false)
		}

		res := cl.Classifier.Load(modelFile)
		if !res {
			slog.Error("classifier load failed")
		}
		return wypes.Bool(true)
	}
}

func detectMultiScaleCascadeClassifierFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*CascadeClassifier], wypes.HostRef[*Frame], wypes.ReturnedList[Rect]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*CascadeClassifier], mat wypes.HostRef[*Frame], list wypes.ReturnedList[Rect]) wypes.Void {
		cl := ref.Raw
		if cl == nil {
			slog.Error("classifier ref is nil")
			return wypes.Void{}
		}
		rects := cl.Classifier.DetectMultiScale(mat.Raw.Image)

		result := make([]Rect, len(rects))
		for i, rect := range rects {
			result[i] = Rect{
				Min: Size{X: wypes.Int32(rect.Min.X), Y: wypes.Int32(rect.Min.Y)},
				Max: Size{X: wypes.Int32(rect.Max.X), Y: wypes.Int32(rect.Max.Y)},
			}
		}

		list.Raw = result
		list.DataPtr = ctx.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func newFaceDetectorYNFunc[T *FaceDetectorYN](ctx *Context) func(*wypes.Store, wypes.String, wypes.String, Size) wypes.HostRef[T] {
	return func(s *wypes.Store, model wypes.String, config wypes.String, sz Size) wypes.HostRef[T] {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			slog.Info(fmt.Sprintf("Downloading model %s...", name))

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				slog.Error(fmt.Sprintf("Error downloading model %v", err))
				return wypes.HostRef[T]{}
			}

		case !models.ModelExists(modelFile):
			return wypes.HostRef[T]{}
		}

		f := gocv.NewFaceDetectorYN(modelFile, "", sz.Unwrap())

		fd := NewFaceDetectorYN(modelFile)
		fd.SetDetector(f)

		v := wypes.HostRef[T]{Raw: fd}
		return v
	}
}

func closeFaceDetectorYNFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*FaceDetectorYN]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*FaceDetectorYN]) wypes.Void {
		f := ref.Raw
		f.Close()

		return wypes.Void{}
	}
}

func faceDetectorYNSetInputSizeFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*FaceDetectorYN], Size) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*FaceDetectorYN], sz Size) wypes.Void {
		f := ref.Raw
		f.Detector.SetInputSize(sz.Unwrap())

		return wypes.Void{}
	}
}

func faceDetectorYNDetectFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*FaceDetectorYN], wypes.HostRef[*Frame]) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*FaceDetectorYN], mat wypes.HostRef[*Frame]) wypes.HostRef[*Frame] {
		f := ref.Raw

		dst := NewEmptyFrame()
		f.Detector.Detect(mat.Unwrap().Image, &dst.Image)

		v := wypes.HostRef[*Frame]{Raw: dst}

		return v
	}
}
