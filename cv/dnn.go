package cv

import (
	"image"
	"log"

	"github.com/wasmvision/wasmvision/models"
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func NetModules(config *Config) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/dnn": wypes.Module{
			"[static]net.read":                       wypes.H3(netReadNetFunc(config)),
			"[static]net.read-from-onnx":             wypes.H2(netReadNetFromONNXFunc(config)),
			"[method]net.close":                      wypes.H2(netCloseFunc(config)),
			"[method]net.empty":                      wypes.H2(netEmptyFunc(config)),
			"[method]net.set-input":                  wypes.H4(netSetInputFunc(config)),
			"[method]net.forward":                    wypes.H3(netForwardFunc(config)),
			"[method]net.get-unconnected-out-layers": wypes.H3(netGetUnconnectedOutLayersFunc(config)),
			"blob-from-image":                        wypes.H11(netBlobFromImageFunc(config)),
		},
	}
}

func netReadNetFunc[T *Net](conf *Config) func(store wypes.Store, model wypes.String, config wypes.String) wypes.HostRef[T] {
	return func(store wypes.Store, model wypes.String, config wypes.String) wypes.HostRef[T] {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, conf.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			if conf.Logging {
				log.Printf("Downloading model %s...\n", name)
			}

			if err := models.Download(name, conf.ModelsDir); err != nil {
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
		id := store.Refs.Put(v)
		nt.ID = wypes.UInt32(id)

		return v
	}
}

func netReadNetFromONNXFunc[T *Net](conf *Config) func(wypes.Store, wypes.String) wypes.HostRef[T] {
	return func(store wypes.Store, model wypes.String) wypes.HostRef[T] {
		name := model.Unwrap()
		modelFile := models.ModelFileName(name, conf.ModelsDir)

		switch {
		case !models.ModelExists(modelFile) && models.ModelWellKnown(name):
			if conf.Logging {
				log.Printf("Downloading model %s...\n", name)
			}

			if err := models.Download(name, conf.ModelsDir); err != nil {
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
		id := store.Refs.Put(v)
		nt.ID = wypes.UInt32(id)

		return v
	}
}

func netCloseFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Net]) wypes.Void {
	return func(store wypes.Store, ref wypes.HostRef[*Net]) wypes.Void {
		nt := ref.Raw
		nt.Close()

		return wypes.Void{}
	}
}

func netEmptyFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Net]) wypes.Bool {
	return func(store wypes.Store, ref wypes.HostRef[*Net]) wypes.Bool {
		nt := ref.Raw
		return wypes.Bool(nt.Net.Empty())
	}
}

func netSetInputFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Net], wypes.HostRef[*Frame], wypes.String) wypes.Void {
	return func(store wypes.Store, ref wypes.HostRef[*Net], blob wypes.HostRef[*Frame], name wypes.String) wypes.Void {
		nt := ref.Raw
		bl := blob.Raw
		blb := bl.Image

		nt.Net.SetInput(blb, name.Unwrap())

		return wypes.Void{}
	}
}

func netForwardFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Net], wypes.String) wypes.HostRef[*Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*Net], output wypes.String) wypes.HostRef[*Frame] {
		nt := ref.Raw

		dst := NewFrame(nt.Net.Forward(output.Unwrap()))

		v := wypes.HostRef[*Frame]{Raw: dst}
		id := store.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		return v
	}
}

func netGetUnconnectedOutLayersFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Net], wypes.List[uint32]) wypes.Void {
	return func(store wypes.Store, ref wypes.HostRef[*Net], list wypes.List[uint32]) wypes.Void {
		nt := ref.Raw

		ls := nt.Net.GetUnconnectedOutLayers()
		result := make([]uint32, len(ls))
		for i, l := range ls {
			result[i] = uint32(l)
		}

		list.Raw = result
		list.DataPtr = conf.ReturnDataPtr
		list.Lower(store)

		return wypes.Void{}
	}
}

func netBlobFromImageFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.Float32, wypes.UInt32, wypes.UInt32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Float32, wypes.Bool, wypes.Bool) wypes.HostRef[*Frame] {
	return func(store wypes.Store, ref wypes.HostRef[*Frame], scale wypes.Float32, size0 wypes.UInt32, size1 wypes.UInt32, mean0 wypes.Float32, mean1 wypes.Float32, mean2 wypes.Float32, mean3 wypes.Float32, swapRb wypes.Bool, crop wypes.Bool) wypes.HostRef[*Frame] {
		frm := ref.Raw

		b := gocv.BlobFromImage(frm.Image, float64(scale.Unwrap()), image.Pt(int(size0.Unwrap()), int(size1.Unwrap())), gocv.NewScalar(float64(mean0.Unwrap()), float64(mean1.Unwrap()), float64(mean2.Unwrap()), float64(mean3.Unwrap())), swapRb.Unwrap(), crop.Unwrap())

		blob := NewFrame(b)

		v := wypes.HostRef[*Frame]{Raw: blob}
		id := store.Refs.Put(v)
		blob.ID = wypes.UInt32(id)

		return v
	}
}
