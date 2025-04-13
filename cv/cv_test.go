package cv

import (
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/config"
)

func getTestStore() wypes.Store {
	stack := wypes.NewSliceStack(64)
	return wypes.Store{
		Stack:  stack,
		Memory: wypes.NewSliceMemory(4096),
		Refs:   wypes.NewMapRefs(),
	}
}

func getTestCVContext() *Context {
	configStore := config.NewStore(map[string]string{})
	ctx := NewContext("test", configStore, "memory", false)
	ctx.ReturnDataPtr = 512
	return ctx
}

func testFrameResult() wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32] {
	return wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]{
		DataPtr: 512,
		Offset:  64,
	}
}

func getResult(store *wypes.Store, res wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32]) wypes.Result[wypes.HostRef[*Frame], wypes.HostRef[*Frame], wypes.UInt32] {
	store.Stack.Push(64)
	return res.Lift(store)
}
