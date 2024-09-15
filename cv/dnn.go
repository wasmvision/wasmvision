package cv

import (
	"github.com/wasmvision/wasmvision/engine"

	"github.com/orsinium-labs/wypes"
)

func NetModules() wypes.Modules {
	return wypes.Modules{
		"wasm:cv/net": wypes.Module{
			"[method]net.close": wypes.H1(netCloseFunc),
			"[method]net.empty": wypes.H1(netEmptyFunc),
		},
	}
}

func netCloseFunc(ref wypes.UInt32) wypes.Void {
	n, ok := engine.NetCache[ref]
	if !ok {
		return wypes.Void{}
	}
	net := n.Net

	net.Close()

	return wypes.Void{}
}

func netEmptyFunc(ref wypes.UInt32) wypes.Bool {
	n, ok := engine.NetCache[ref]
	if !ok {
		return wypes.Bool(true)
	}
	net := n.Net

	return wypes.Bool(net.Empty())
}
