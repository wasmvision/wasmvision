package cv

import (
	"gocv.io/x/gocv"

	"github.com/orsinium-labs/wypes"
)

func MatModules(config *Config) wypes.Modules {
	return wypes.Modules{
		"wasm:cv/mat": wypes.Module{
			"[constructor]mat":          wypes.H4(matNewWithSizeFunc(config)),
			"[resource-drop]mat":        wypes.H2(matCloseFunc(config)),
			"[static]mat.new-mat":       wypes.H1(matNewFunc(config)),
			"[static]mat.new-with-size": wypes.H4(matNewWithSizeFunc(config)),
			"[method]mat.close":         wypes.H2(matCloseFunc(config)),
			"[method]mat.cols":          wypes.H2(matColsFunc(config)),
			"[method]mat.rows":          wypes.H2(matRowsFunc(config)),
			"[method]mat.mattype":       wypes.H2(matTypeFunc(config)),
			"[method]mat.empty":         wypes.H2(matEmptyFunc(config)),
			"[method]mat.size":          wypes.H3(matSizeFunc(config)),
			"[method]mat.reshape":       wypes.H4(matReshapeFunc(config)),
			"[method]mat.get-float-at":  wypes.H4(matGetFloatAtFunc(config)),
			"[method]mat.set-float-at":  wypes.H5(matSetFloatAtFunc(config)),
			"[method]mat.get-uchar-at":  wypes.H4(matGetUcharAtFunc(config)),
			"[method]mat.set-uchar-at":  wypes.H5(matSetUcharAtFunc(config)),
			"[method]mat.get-vecb-at":   wypes.H5(matGetVecbAtFunc(config)),
		},
	}
}

func matNewFunc(conf *Config) func(wypes.Store) wypes.HostRef[*Frame] {
	return func(s wypes.Store) wypes.HostRef[*Frame] {
		f := NewEmptyFrame()

		v := wypes.HostRef[*Frame]{Raw: f}
		id := s.Refs.Put(v)
		f.ID = wypes.UInt32(id)

		return v
	}
}

func matNewWithSizeFunc(conf *Config) func(wypes.Store, wypes.UInt32, wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s wypes.Store, rows, cols, matType wypes.UInt32) wypes.HostRef[*Frame] {
		f := NewFrame(gocv.NewMatWithSize(int(rows), int(cols), gocv.MatType(matType)))

		v := wypes.HostRef[*Frame]{Raw: f}
		id := s.Refs.Put(v)
		f.ID = wypes.UInt32(id)

		return v
	}
}

func matCloseFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame]) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.Void {
		f := ref.Raw
		f.Close()

		return wypes.Void{}
	}
}

func matColsFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Cols())
	}
}

func matRowsFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Rows())
	}
}

func matTypeFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame]) wypes.UInt32 {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.UInt32 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt32(mat.Type())
	}
}

func matEmptyFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame]) wypes.Bool {
	return func(s wypes.Store, ref wypes.HostRef[*Frame]) wypes.Bool {
		f := ref.Raw
		mat := f.Image

		return wypes.Bool(mat.Empty())
	}
}

func matSizeFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.List[uint32]) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], list wypes.List[uint32]) wypes.Void {
		f := ref.Raw
		mat := f.Image
		dims := mat.Size()

		result := make([]uint32, len(dims))
		for i, dim := range dims {
			result[i] = uint32(dim)
		}

		list.Raw = result
		list.DataPtr = conf.ReturnDataPtr
		list.Lower(s)

		return wypes.Void{}
	}
}

func matReshapeFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.HostRef[*Frame] {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], channels wypes.UInt32, rows wypes.UInt32) wypes.HostRef[*Frame] {
		f := ref.Raw
		mat := f.Image

		d := mat.Reshape(int(channels), int(rows))
		dst := NewFrame(d)

		v := wypes.HostRef[*Frame]{Raw: dst}
		id := s.Refs.Put(v)
		dst.ID = wypes.UInt32(id)

		return v
	}
}

func matGetFloatAtFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.Float32 {
	return func(store wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32) wypes.Float32 {
		f := ref.Raw
		mat := f.Image

		return wypes.Float32(mat.GetFloatAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetFloatAtFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.Float32) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.Float32) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.SetFloatAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetUcharAtFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32) wypes.UInt8 {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32) wypes.UInt8 {
		f := ref.Raw
		mat := f.Image

		return wypes.UInt8(mat.GetUCharAt(int(row.Unwrap()), int(col.Unwrap())))
	}
}

func matSetUcharAtFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.UInt8) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.UInt8) wypes.Void {
		f := ref.Raw
		mat := f.Image
		mat.SetUCharAt(int(row.Unwrap()), int(col.Unwrap()), v.Unwrap())

		return wypes.Void{}
	}
}

func matGetVecbAtFunc(conf *Config) func(wypes.Store, wypes.HostRef[*Frame], wypes.UInt32, wypes.UInt32, wypes.List[uint8]) wypes.Void {
	return func(s wypes.Store, ref wypes.HostRef[*Frame], row wypes.UInt32, col wypes.UInt32, v wypes.List[uint8]) wypes.Void {
		f := ref.Raw
		mat := f.Image
		data := mat.GetVecbAt(int(row.Unwrap()), int(col.Unwrap()))

		result := make([]uint8, len(data))
		for i, dim := range data {
			result[i] = uint8(dim)
		}

		v.Raw = result
		v.DataPtr = conf.ReturnDataPtr
		v.Lower(s)

		return wypes.Void{}
	}
}
