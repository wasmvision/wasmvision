package runtime

import (
	"fmt"
	"log/slog"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
)

// hostFramedataOpenFunc is just here to fulfill the interface for the host datastore constructor function.
func hostFramedataOpenFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32) wypes.UInt32 {
	return func(*wypes.Store, wypes.UInt32) wypes.UInt32 {
		return wypes.UInt32(1)
	}
}

// hostFramedataDropFunc is just here to fulfill the interface for the host datastore drop function.
func hostFramedataDropFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32) wypes.Void {
	return func(*wypes.Store, wypes.UInt32) wypes.Void {
		return wypes.Void{}
	}
}

// hostFramedataDeleteFunc deletes a key from the frame store for a given frame.
func hostFramedataDeleteFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.String, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, frame wypes.UInt32, key wypes.String, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		ctx.FrameStore.Delete(int(frame.Unwrap()), key.Unwrap())

		result.IsError = false
		result.OK = 0

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostFramedataDeleteFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostFramedataExistsFunc checks if there is any data in the frame store for a given frame.
func hostFramedataExistsFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.Result[wypes.Bool, wypes.Bool, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, frame wypes.UInt32, result wypes.Result[wypes.Bool, wypes.Bool, wypes.UInt32]) wypes.Void {
		ok := ctx.FrameStore.Exists(int(frame.Unwrap()))
		result.IsError = false
		result.OK = wypes.Bool(ok)

		return wypes.Void{}
	}
}

// hostFramedataGetFunc gets the data for a key from the frame store for a given frame.
func hostFramedataGetFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.String, wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, frame wypes.UInt32, key wypes.String, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		val, ok := ctx.FrameStore.Get(int(frame.Unwrap()), key.Unwrap())

		if !ok {
			result.IsError = true
			result.Error = wypes.UInt32(2) // no-such-key
		} else {
			result.IsError = false
			result.OK = wypes.Bytes{Raw: []byte(val)}
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostFramedataGetFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostFramedataGetKeysFunc gets all the keys for a given frame from the frame store.
func hostFramedataGetKeysFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.ReturnedList[wypes.String]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, frame wypes.UInt32, result wypes.ReturnedList[wypes.String]) wypes.Void {
		keys, ok := ctx.FrameStore.GetKeys(int(frame.Unwrap()))

		if !ok {
			// no data for this frame
			result.Raw = []wypes.String{}
		} else {
			res := make([]wypes.String, 0, len(keys))
			for i := range keys {
				v := keys[i]
				res = append(res, wypes.String{Raw: v})
			}
			result.Raw = res
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostFramedataGetKeysFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostFramedataSetFunc sets the data for a key in the frame store for a given frame.
func hostFramedataSetFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.UInt32, wypes.String, wypes.Bytes, wypes.Result[wypes.UInt32, wypes.Bool, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, frame wypes.UInt32, key wypes.String, data wypes.Bytes, result wypes.Result[wypes.UInt32, wypes.Bool, wypes.UInt32]) wypes.Void {
		err := ctx.FrameStore.Set(int(frame.Unwrap()), key.Unwrap(), string(data.Raw))

		if err != nil {
			slog.Error(fmt.Sprintf("hostFramedataSetFunc error in store after set: %v", err))
			result.IsError = true
			result.Error = wypes.UInt32(1)
		} else {
			result.IsError = false
			result.OK = wypes.Bool(true)
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostFramedataSetFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostProcessorOpenFunc is just here to fulfill the interface for the host datastore constructor function.
func hostProcessorDataOpenFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32) wypes.UInt32 {
	return func(*wypes.Store, wypes.UInt32) wypes.UInt32 {
		return wypes.UInt32(1)
	}
}

// hostProcessorDataDropFunc is just here to fulfill the interface for the host datastore drop function.
func hostProcessorDataDropFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32) wypes.Void {
	return func(*wypes.Store, wypes.UInt32) wypes.Void {
		return wypes.Void{}
	}
}

// hostProcessorDataDeleteFunc deletes a key from the processor store for a given processor.
func hostProcessorDataDeleteFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.String, wypes.String, wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, processor wypes.String, key wypes.String, result wypes.Result[wypes.UInt32, wypes.UInt32, wypes.UInt32]) wypes.Void {
		ctx.ProcessorStore.Delete(processor.Unwrap(), key.Unwrap())

		result.IsError = false
		result.OK = 0

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostProcessorDataDeleteFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostProcessorDataExistsFunc checks if there is any data in the processor store for a given processor.
func hostProcessorDataExistsFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.String, wypes.Result[wypes.Bool, wypes.Bool, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, processor wypes.String, result wypes.Result[wypes.Bool, wypes.Bool, wypes.UInt32]) wypes.Void {
		ok := ctx.ProcessorStore.Exists(processor.Unwrap())
		result.IsError = false
		result.OK = wypes.Bool(ok)

		return wypes.Void{}
	}
}

// hostProcessorDataGetFunc gets the data for a key from the processor store for a given processor.
func hostProcessorDataGetFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.String, wypes.String, wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, processor wypes.String, key wypes.String, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		val, ok := ctx.ProcessorStore.Get(processor.Unwrap(), key.Unwrap())

		if !ok {
			result.IsError = true
			result.Error = wypes.UInt32(2) // no-such-key
		} else {
			result.IsError = false
			result.OK = wypes.Bytes{Raw: val}
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostProcessorDataGetFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostProcessorDataGetKeysFunc gets all the keys for a given processor from the processor store.
func hostProcessorDataGetKeysFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.String, wypes.ReturnedList[wypes.String]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, processor wypes.String, result wypes.ReturnedList[wypes.String]) wypes.Void {
		keys, ok := ctx.ProcessorStore.GetKeys(processor.Unwrap())

		if !ok {
			// no data for this processor
			result.Raw = []wypes.String{}
		} else {
			res := make([]wypes.String, 0, len(keys))
			for i := range keys {
				v := keys[i]
				res = append(res, wypes.String{Raw: v})
			}
			result.Raw = res
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostProcessorDataGetKeysFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

// hostProcessorDataSetFunc sets the data for a key in the processor store for a given processor.
func hostProcessorDataSetFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32, wypes.String, wypes.String, wypes.Bytes, wypes.Result[wypes.UInt32, wypes.Bool, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, fs wypes.UInt32, processor wypes.String, key wypes.String, data wypes.Bytes, result wypes.Result[wypes.UInt32, wypes.Bool, wypes.UInt32]) wypes.Void {
		value := make([]byte, 0, len(data.Raw))
		value = append(value, data.Raw...)
		err := ctx.ProcessorStore.Set(processor.Unwrap(), key.Unwrap(), value)
		if err != nil {
			slog.Error(fmt.Sprintf("hostProcessorDataSetFunc error in store after set: %v", err))
			result.IsError = true
			result.Error = wypes.UInt32(1)
		} else {
			result.IsError = false
			result.OK = wypes.Bool(true)
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("hostProcessorDataSetFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}
