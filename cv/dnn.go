package cv

import (
	"image"
	"log"

	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func NetModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/dnn": wypes.Module{
			"[static]net.read":                       wypes.H3(netReadNetFunc(ctx)),
			"[static]net.read-from-onnx":             wypes.H2(netReadNetFromONNXFunc(ctx)),
			"[method]net.close":                      wypes.H2(netCloseFunc(ctx)),
			"[method]net.empty":                      wypes.H2(netEmptyFunc(ctx)),
			"[method]net.set-input":                  wypes.H4(netSetInputFunc(ctx)),
			"[method]net.forward":                    wypes.H3(netForwardFunc(ctx)),
			"[method]net.get-unconnected-out-layers": wypes.H3(netGetUnconnectedOutLayersFunc(ctx)),
			"blob-from-image":                        wypes.H11(netBlobFromImageFunc(ctx)),
		},
	}
}

func netReadNetFunc[T *Net](ctx *Context) func(*wypes.Store, wypes.String, wypes.String) wypes.HostRef[T] {
	return func(s *wypes.Store, model wypes.String, config wypes.String) wypes.HostRef[T] {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			if ctx.Logging {
				log.Printf("Downloading model %s...\n", name)
			}

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				log.Printf("Error downloading model: %s", err)
				return wypes.HostRef[T]{}
			}

		case !models.ModelExists(modelFile):
			return wypes.HostRef[T]{}
		}

		n := gocv.ReadNet(modelFile, config.Unwrap())
		if n.Empty() {
			return wypes.HostRef[T]{}
		}

		nt := NewNet(model.Unwrap())
		nt.SetNet(n)

		v := wypes.HostRef[T]{Raw: nt}

		return v
	}
}

func netReadNetFromONNXFunc[T *Net](ctx *Context) func(*wypes.Store, wypes.String) wypes.HostRef[T] {
	return func(s *wypes.Store, model wypes.String) wypes.HostRef[T] {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			if ctx.Logging {
				log.Printf("Downloading model %s...\n", name)
			}

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				log.Printf("Error downloading model: %s", err)
				return wypes.HostRef[T]{}
			}

		case !models.ModelExists(modelFile):
			return wypes.HostRef[T]{}
		}

		n := gocv.ReadNetFromONNX(modelFile)
		if n.Empty() {
			return wypes.HostRef[T]{}
		}

		nt := NewNet(model.Unwrap())
		nt.SetNet(n)

		v := wypes.HostRef[T]{Raw: nt}

		return v
	}
}

func netCloseFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}

func netEmptyFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net]) wypes.Bool {
	return func(s *wypes.Store, ref wypes.HostRef[*Net]) wypes.Bool {
		nt := ref.Raw
		return wypes.Bool(nt.Net.Empty())
	}
}

func netSetInputFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.HostRef[*Frame], wypes.String) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], blob wypes.HostRef[*Frame], name wypes.String) wypes.Void {
		nt := ref.Raw
		bl := blob.Raw
		blb := bl.Image

		nt.Net.SetInput(blb, name.Unwrap())

		return wypes.Void{}
	}
}

func netForwardFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.String) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], output wypes.String) wypes.HostRef[*Frame] {
		nt := ref.Raw

		dst := NewFrame(nt.Net.Forward(output.Unwrap()))

		v := wypes.HostRef[*Frame]{Raw: dst}

		return v
	}
}

func netGetUnconnectedOutLayersFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.ReturnedList[wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], list wypes.ReturnedList[wypes.UInt32]) wypes.Void {
		nt := ref.Raw

		ls := nt.Net.GetUnconnectedOutLayers()
		result := make([]wypes.UInt32, len(ls))
		for i, l := range ls {
			result[i] = wypes.UInt32(l)
		}

		list.Raw = result
		list.DataPtr = ctx.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func netBlobFromImageFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Bool, wypes.Bool) wypes.HostRef[*Frame] {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], scale wypes.Float32, size0 wypes.UInt32, size1 wypes.UInt32, mean0 wypes.Float32, mean1 wypes.Float32, mean2 wypes.Float32, mean3 wypes.Float32, swapRb wypes.Bool, crop wypes.Bool) wypes.HostRef[*Frame] {
		frm := ref.Raw

		b := gocv.BlobFromImage(frm.Image, float64(scale.Unwrap()), image.Pt(int(size0.Unwrap()), int(size1.Unwrap())), gocv.NewScalar(float64(mean0.Unwrap()), float64(mean1.Unwrap()), float64(mean2.Unwrap()), float64(mean3.Unwrap())), swapRb.Unwrap(), crop.Unwrap())

		blob := NewFrame(b)

		v := wypes.HostRef[*Frame]{Raw: blob}

		return v
	}
}
