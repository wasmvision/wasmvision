package cv

import (
	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/frame"
)

func ObjDetectModules(config *Config) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/objdetect": wypes.Module{
			"[constructor]cascade-classifier":  wypes.H1(newCascadeClassifierFunc(config)),
			"[method]cascade-classifier.close": wypes.H2(closeCascadeClassifierFunc(config)),
		},
	}
}

func newCascadeClassifierFunc(conf *Config) func(wypes.Store) wypes.HostRef[*frame.Frame] {
	return func(s wypes.Store) wypes.HostRef[*frame.Frame] {
		f := frame.NewEmptyFrame()

		v := wypes.HostRef[*frame.Frame]{Raw: f}
		id := s.Refs.Put(v)
		f.ID = wypes.UInt32(id)

		return v
	}
}

func closeCascadeClassifierFunc(conf *Config) func(wypes.Store, wypes.HostRef[*frame.Frame]) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*frame.Frame]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}
