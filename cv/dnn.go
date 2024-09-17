package cv

import (
	"github.com/wasmvision/wasmvision/frame"
	"github.com/wasmvision/wasmvision/net"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func NetModules(fc *frame.Cache, nc *net.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/net": wypes.Module{
			"[static]net.read-net-from-onnx": wypes.H1(netReadNetFromONNXFunc(nc)),
			"[method]net.close":              wypes.H1(netCloseFunc(nc)),
			"[method]net.empty":              wypes.H1(netEmptyFunc(nc)),
			"[method]net.set-input":          wypes.H3(netSetInputFunc(nc, fc)),
			"[method]net.forward":            wypes.H2(netForwardFunc(nc, fc)),
		},
	}
}

func netReadNetFromONNXFunc(cache *net.Cache) func(onnxModelPath wypes.String) wypes.UInt32 {
	return func(onnxModelPath wypes.String) wypes.UInt32 {
		n := gocv.ReadNetFromONNX(onnxModelPath.Unwrap())
		if n.Empty() {
			return wypes.UInt32(0)
		}
		net := net.NewNet()
		net.SetNet(n)

		cache.Set(net)

		return wypes.UInt32(net.ID)
	}
}

func netCloseFunc(cache *net.Cache) func(ref wypes.UInt32) wypes.Void {
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

func netEmptyFunc(cache *net.Cache) func(ref wypes.UInt32) wypes.Bool {
	return func(ref wypes.UInt32) wypes.Bool {
		n, ok := cache.Get(ref)
		if !ok {
			return wypes.Bool(true)
		}
		net := n.Net

		return wypes.Bool(net.Empty())
	}
}

func netSetInputFunc(cache *net.Cache, framecache *frame.Cache) func(ref wypes.UInt32, blob wypes.UInt32, name wypes.String) wypes.Void {
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

func netForwardFunc(cache *net.Cache, framecache *frame.Cache) func(ref wypes.UInt32, output wypes.String) wypes.UInt32 {
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
