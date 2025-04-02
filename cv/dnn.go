package cv

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func NetModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/dnn": wypes.Module{
			"[static]net.read":                       wypes.H4(netReadNetFunc(ctx)),
			"[static]net.read-from-onnx":             wypes.H3(netReadNetFromONNXFunc(ctx)),
			"[resource-drop]net":                     wypes.H2(netCloseFunc(ctx)),
			"[method]net.close":                      wypes.H2(netCloseFunc(ctx)),
			"[method]net.empty":                      wypes.H2(netEmptyFunc(ctx)),
			"[method]net.set-input":                  wypes.H5(netSetInputFunc(ctx)),
			"[method]net.forward":                    wypes.H4(netForwardFunc(ctx)),
			"[method]net.get-unconnected-out-layers": wypes.H3(netGetUnconnectedOutLayersFunc(ctx)),
			"blob-from-image":                        wypes.H8(netBlobFromImageFunc(ctx)),
			"[method]net.get-layer":                  wypes.H4(netGetLayerFunc(ctx)),
			"[method]layer.get-name":                 wypes.H3(netLayerGetNameFunc(ctx)),
			//"[method]layer.close":                    wypes.H2(netLayerCloseFunc(ctx)),
		},
	}
}

func netReadNetFunc[T *Net](ctx *Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, model wypes.String, config wypes.String, result wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32]) wypes.Void {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			slog.Info(fmt.Sprintf("Downloading model %s...", name))

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				handleNetError(ctx, s, nil, result, err)
				return wypes.Void{}
			}

		case !models.ModelExists(modelFile):
			handleNetError(ctx, s, nil, result, fmt.Errorf("model %s not found", name))
			return wypes.Void{}
		}

		n := gocv.ReadNet(modelFile, config.Unwrap())
		if n.Empty() {
			handleNetError(ctx, s, nil, result, fmt.Errorf("model %s not read", name))
			return wypes.Void{}
		}

		backend, target := gocv.NetBackendDefault, gocv.NetTargetCPU
		if ctx.EnableCUDA {
			backend = gocv.NetBackendCUDA
			target = gocv.NetTargetCUDA
		}
		n.SetPreferableBackend(backend)
		n.SetPreferableTarget(target)

		nt := NewNet(model.Unwrap())
		nt.SetNet(n)

		handleNetReturn(ctx, s, nt, result)
		return wypes.Void{}
	}
}

func netReadNetFromONNXFunc[T *Net](ctx *Context) func(*wypes.Store, wypes.String, wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, model wypes.String, result wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32]) wypes.Void {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, ctx.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			slog.Info(fmt.Sprintf("Downloading model %s...", name))

			if err := models.Download(name, ctx.ModelsDir); err != nil {
				handleNetError(ctx, s, nil, result, err)
				return wypes.Void{}
			}

		case !models.ModelExists(modelFile):
			handleNetError(ctx, s, nil, result, fmt.Errorf("model %s not found", name))
			return wypes.Void{}
		}

		n := gocv.ReadNetFromONNX(modelFile)
		if n.Empty() {
			handleNetError(ctx, s, nil, result, fmt.Errorf("model %s not read", name))
			return wypes.Void{}
		}

		backend, target := gocv.NetBackendDefault, gocv.NetTargetCPU
		if ctx.EnableCUDA {
			backend = gocv.NetBackendCUDA
			target = gocv.NetTargetCUDA
		}
		n.SetPreferableBackend(backend)
		n.SetPreferableTarget(target)

		nt := NewNet(model.Unwrap())
		nt.SetNet(n)

		handleNetReturn(ctx, s, nt, result)
		return wypes.Void{}
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

func netSetInputFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.HostRef[*Frame], wypes.String, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], blob wypes.HostRef[*Frame], name wypes.String, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		nt := ref.Raw
		bl := blob.Raw
		blb := bl.Image

		nt.Net.SetInput(blb, name.Unwrap())

		return wypes.Void{}
	}
}

func netForwardFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.String, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], output wypes.String, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		nt := ref.Raw
		dst := NewFrame(nt.Net.Forward(output.Unwrap()))
		if dst.Image.Empty() {
			handleFrameError(ctx, s, dst, result, errors.New("empty forward output"))
			return wypes.Void{}
		}

		handleFrameReturn(ctx, s, dst, result)
		return wypes.Void{}
	}
}

func netGetUnconnectedOutLayersFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.Result[wypes.List[wypes.UInt32], wypes.List[wypes.UInt32], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], result wypes.Result[wypes.List[wypes.UInt32], wypes.List[wypes.UInt32], wypes.UInt32]) wypes.Void {
		nt := ref.Raw

		ls := nt.Net.GetUnconnectedOutLayers()
		list := make([]wypes.UInt32, len(ls))
		for i, l := range ls {
			list[i] = wypes.UInt32(l)
		}

		result.IsError = false
		result.OK = wypes.List[wypes.UInt32]{Raw: list}
		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		return wypes.Void{}
	}
}

func netBlobFromImageFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Frame], wypes.Float32, Size, Scalar, wypes.Bool, wypes.Bool, wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Frame], scale wypes.Float32, sz Size, scalar Scalar, swapRb wypes.Bool, crop wypes.Bool, result wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Void {
		frm := ref.Raw

		b := gocv.BlobFromImage(frm.Image, float64(scale.Unwrap()), sz.Unwrap(), scalar.Unwrap(), swapRb.Unwrap(), crop.Unwrap())

		handleFrameReturn(ctx, s, NewFrame(b), result)
		return wypes.Void{}
	}
}

func netGetLayerFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Net], wypes.UInt32, wypes.Result[wypes.HostRef[*Layer], wypes.HostRef[*Layer], wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Net], id wypes.UInt32, result wypes.Result[wypes.HostRef[*Layer], wypes.HostRef[*Layer], wypes.UInt32]) wypes.Void {
		nt := ref.Raw

		dst := NewLayer(nt.Net.GetLayer(int(id.Unwrap())))

		result.IsError = false
		result.OK = wypes.HostRef[*Layer]{Raw: dst}
		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		return wypes.Void{}
	}
}

func netLayerGetNameFunc(ctx *Context) func(*wypes.Store, wypes.HostRef[*Layer], wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, ref wypes.HostRef[*Layer], result wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
		layer := ref.Raw

		name := layer.Layer.GetName()

		result.IsError = false
		result.OK = wypes.String{Raw: name}
		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		return wypes.Void{}
	}
}
