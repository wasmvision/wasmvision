package cv

import (
	"image"

	"github.com/wasmvision/wasmvision/frame"
	"github.com/wasmvision/wasmvision/net"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func NetModules(fc *frame.Cache, nc *net.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/dnn": wypes.Module{
			"[static]net.read-net":                   wypes.H2(netReadNetFunc(nc)),
			"[static]net.read-net-from-onnx":         wypes.H1(netReadNetFromONNXFunc(nc)),
			"[method]net.close":                      wypes.H1(netCloseFunc(nc)),
			"[method]net.empty":                      wypes.H1(netEmptyFunc(nc)),
			"[method]net.set-input":                  wypes.H3(netSetInputFunc(nc, fc)),
			"[method]net.forward":                    wypes.H2(netForwardFunc(nc, fc)),
			"[method]net.get-unconnected-out-layers": wypes.H3(netGetUnconnectedOutLayersFunc(nc, fc)),
			"blob-from-image":                        wypes.H10(netBlobFromImageFunc(fc)),
		},
	}
}

func netReadNetFunc(cache *net.Cache) func(wypes.String, wypes.String) wypes.UInt32 {
	return func(model wypes.String, config wypes.String) wypes.UInt32 {
		modelFile := cache.ModelFileName(model.Unwrap())

		n := gocv.ReadNet(modelFile, config.Unwrap())
		if n.Empty() {
			return wypes.UInt32(0)
		}

		net := net.NewNet(model.Unwrap())
		net.SetNet(n)
		cache.Set(net)

		return wypes.UInt32(net.ID)
	}
}

func netReadNetFromONNXFunc(cache *net.Cache) func(wypes.String) wypes.UInt32 {
	return func(model wypes.String) wypes.UInt32 {
		modelFile := cache.ModelFileName(model.Unwrap())

		n := gocv.ReadNetFromONNX(modelFile)
		if n.Empty() {
			return wypes.UInt32(0)
		}

		net := net.NewNet(model.Unwrap())
		net.SetNet(n)
		cache.Set(net)

		return wypes.UInt32(net.ID)
	}
}

func netCloseFunc(cache *net.Cache) func(wypes.UInt32) wypes.Void {
	return func(ref wypes.UInt32) wypes.Void {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}

		net := n.Net
		net.Close()

		return wypes.Void{}
	}
}

func netEmptyFunc(cache *net.Cache) func(wypes.UInt32) wypes.Bool {
	return func(ref wypes.UInt32) wypes.Bool {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.Bool(true)
		}
		net := n.Net

		return wypes.Bool(net.Empty())
	}
}

func netSetInputFunc(cache *net.Cache, framecache *frame.Cache) func(wypes.UInt32, wypes.UInt32, wypes.String) wypes.Void {
	return func(ref wypes.UInt32, blob wypes.UInt32, name wypes.String) wypes.Void {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		nt := n.Net

		b, ok := framecache.Get(blob)
		if !ok {
			return wypes.Void{}
		}
		blb := b.Image

		nt.SetInput(blb, name.Unwrap())

		return wypes.Void{}
	}
}

func netForwardFunc(cache *net.Cache, framecache *frame.Cache) func(wypes.UInt32, wypes.String) wypes.UInt32 {
	return func(ref wypes.UInt32, output wypes.String) wypes.UInt32 {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.UInt32(0)
		}
		nt := n.Net

		dst := frame.NewFrame()
		dst.SetImage(nt.Forward(output.Unwrap()))
		framecache.Set(dst)

		return wypes.UInt32(dst.ID)
	}
}

func netGetUnconnectedOutLayersFunc(cache *net.Cache, framecache *frame.Cache) func(wypes.Store, wypes.UInt32, wypes.List[uint32]) wypes.Void {
	return func(s wypes.Store, ref wypes.UInt32, list wypes.List[uint32]) wypes.Void {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.Void{}
		}
		nt := n.Net

		ls := nt.GetUnconnectedOutLayers()
		result := make([]uint32, len(ls))
		for i, l := range ls {
			result[i] = uint32(l)
		}

		list.Raw = result
		list.DataPtr = framecache.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func netBlobFromImageFunc(cache *frame.Cache) func(wypes.UInt32, wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Bool, wypes.Bool) wypes.UInt32 {
	return func(matref wypes.UInt32, scale wypes.Float32, size0 wypes.UInt32, size1 wypes.UInt32, mean0 wypes.Float32, mean1 wypes.Float32, mean2 wypes.Float32, mean3 wypes.Float32, swapRb wypes.Bool, crop wypes.Bool) wypes.UInt32 {
		f, ok := cache.Get(matref)
		if !ok {
			return wypes.UInt32(0)
		}
		src := f.Image

		b := gocv.BlobFromImage(src, float64(scale.Unwrap()), image.Pt(int(size0.Unwrap()), int(size1.Unwrap())), gocv.NewScalar(float64(mean0.Unwrap()), float64(mean1.Unwrap()), float64(mean2.Unwrap()), float64(mean3.Unwrap())), swapRb.Unwrap(), crop.Unwrap())

		blob := frame.NewFrame()
		blob.SetImage(b)
		cache.Set(blob)

		return wypes.UInt32(blob.ID)
	}
}
