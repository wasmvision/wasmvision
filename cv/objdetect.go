package cv

import (
	"log"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"
)

func ObjDetectModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/objdetect": wypes.Module{
			"[constructor]cascade-classifier":               wypes.H2(newCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.close":              wypes.H2(closeCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.load":               wypes.H3(loadCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.detect-multi-scale": wypes.H4(detectMultiScaleCascadeClassifierFunc(ctx)),
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
			if ctx.Logging {
				log.Printf("Downloading classifier %s...\n", name)
			}

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				log.Printf("Error downloading classifier: %s", err)
				return wypes.Bool(false)
			}

		case !models.ModelExists(modelFile):
			return wypes.Bool(false)
		}

		if cl == nil {
			log.Println("classifier is nil")
			return wypes.Bool(false)
		}

		res := cl.Classifier.Load(modelFile)
		if !res {
			log.Println("classifier load failed")
		}
		return wypes.Bool(true)
	}
}

func detectMultiScaleCascadeClassifierFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*CascadeClassifier], wypes.HostRef[*Frame], wypes.ReturnedList[Rect]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*CascadeClassifier], mat wypes.HostRef[*Frame], list wypes.ReturnedList[Rect]) wypes.Void {
		cl := ref.Raw
		if cl == nil {
			log.Println("classifier ref is nil", ref)
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
