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
			result.Error = wypes.UInt32(1) // no-such-key
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
