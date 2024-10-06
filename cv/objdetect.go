package cv

import (
	"github.com/orsinium-labs/wypes"
)

func ObjDetectModules(ctx *Context) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/objdetect": wypes.Module{
			"[constructor]cascade-classifier":  wypes.H1(newCascadeClassifierFunc(ctx)),
			"[method]cascade-classifier.close": wypes.H2(closeCascadeClassifierFunc(ctx)),
		},
	}
}

func newCascadeClassifierFunc(ctx *Context) func(wypes.Store) wypes.HostRef[*Frame] {
	return func(s wypes.Store) wypes.HostRef[*Frame] {
		f := NewEmptyFrame()

		v := wypes.HostRef[*Frame]{Raw: f}
		id := s.Refs.Put(v)
		f.ID = wypes.UInt32(id)

		return v
	}
}

func closeCascadeClassifierFunc(ctx *Context) func(wypes.Store, wypes.HostRef[*Frame]) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}
