package cv

import (
	"github.com/wasmvision/wasmvision/net"

	"github.com/orsinium-labs/wypes"
)

func NetModules(cache *net.Cache) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/net": wypes.Module{
			"[method]net.close": wypes.H1(netCloseFunc(cache)),
			"[method]net.empty": wypes.H1(netEmptyFunc(cache)),
		},
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
