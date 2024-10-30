package runtime

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/orsinium-labs/wypes"
	"github.com/wasmvision/wasmvision/cv"
	"gocv.io/x/gocv"
)

func hostedModules(ctx *cv.Context) wypes.Modules {
	return wypes.Modules{
		"wasmvision:platform/config": wypes.Module{
			"get-config": wypes.H3(hostGetConfigFunc(ctx)),
		},
		"wasmvision:platform/http": wypes.Module{
			"get":        wypes.H3(httpGetFunc(ctx)),
			"post":       wypes.H5(httpPostFunc(ctx)),
			"post-image": wypes.H7(httpPostImageFunc(ctx)),
		},
		"wasmvision:platform/logging": wypes.Module{
			"println": wypes.H1(hostPrintln),
			"log":     wypes.H1(hostLogFunc(ctx)),
			"error":   wypes.H1(hostErrorFunc(ctx)),
			"warn":    wypes.H1(hostWarnFunc(ctx)),
			"info":    wypes.H1(hostInfoFunc(ctx)),
			"debug":   wypes.H1(hostDebugFunc(ctx)),
		},
		"wasmvision:platform/time": wypes.Module{
			"now": wypes.H2(hostTimeNowFunc(ctx)),
		},
	}
}

func hostPrintln(msg wypes.String) wypes.Void {
	fmt.Println(msg.Unwrap())
	return wypes.Void{}
}

func hostLogFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		slog.Warn(msg.Unwrap())
		return wypes.Void{}
	}
}

func hostErrorFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		slog.Error(msg.Unwrap())
		return wypes.Void{}
	}
}

func hostWarnFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		slog.Warn(msg.Unwrap())
		return wypes.Void{}
	}
}

func hostInfoFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		slog.Info(msg.Unwrap())
		return wypes.Void{}
	}
}

func hostDebugFunc(ctx *cv.Context) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		slog.Debug(msg.Unwrap())
		return wypes.Void{}
	}
}

func hostGetConfigFunc(ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, key wypes.String, result wypes.Result[wypes.String, wypes.String, wypes.UInt32]) wypes.Void {
		if s.Error != nil {
			slog.Error(fmt.Sprintf("error in store after lift: %v", s.Error))
		}
		val, ok := ctx.Config.Get(key.Unwrap())
		if !ok {
			result.IsError = true
			result.Error = wypes.UInt32(1) // no-such-key
		} else {
			result.IsError = false
			result.OK = wypes.String{Raw: val}
		}

		result.DataPtr = ctx.ReturnDataPtr
		result.Lower(s)

		if s.Error != nil {
			slog.Error(fmt.Sprintf("error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

func httpGetFunc(ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, url wypes.String, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		slog.Info(fmt.Sprintf("http get: %s", url.Unwrap()))

		resp, err := http.Get(url.Unwrap())
		if err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		res := make([]byte, len(body))
		copy(res, body)

		result.IsError = false
		result.OK = wypes.Bytes{Raw: res}
		result.DataPtr = ctx.ReturnDataPtr

		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("httpGetFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

func httpPostFunc(ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Bytes, wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, url wypes.String, contentType wypes.String, request wypes.Bytes, result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		slog.Info(fmt.Sprintf("http post: %s\n", url.Unwrap()))

		reqBody := bytes.NewBuffer(request.Raw)

		resp, err := http.Post(url.Unwrap(), contentType.Unwrap(), reqBody)
		if err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		res := make([]byte, len(body))
		copy(res, body)

		result.IsError = false
		result.OK = wypes.Bytes{Raw: res}
		result.DataPtr = ctx.ReturnDataPtr

		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("httpPostFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

func httpPostImageFunc(ctx *cv.Context) func(*wypes.Store, wypes.String, wypes.String, wypes.Bytes, wypes.String, wypes.HostRef[*cv.Frame], wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
	return func(s *wypes.Store, url wypes.String, contentType wypes.String, template wypes.Bytes, responseKey wypes.String, mat wypes.HostRef[*cv.Frame], result wypes.Result[wypes.Bytes, wypes.Bytes, wypes.UInt32]) wypes.Void {
		slog.Info(fmt.Sprintf("http post image: %s\n", url.Unwrap()))

		// TODO: support other content types
		ct := "application/json"

		buffer, err := gocv.IMEncode(gocv.JPEGFileExt, mat.Raw.Image)
		if err != nil {
			slog.Error(fmt.Sprintf("error encoding image: %v", err))
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		sEnc := base64.StdEncoding.EncodeToString(buffer.GetBytes())

		tmpl := string(template.Raw)
		tmpl = strings.Replace(tmpl, "%IMAGE%", sEnc, -1)

		reqBody := bytes.NewBuffer([]byte(tmpl))

		resp, err := http.Post(url.Unwrap(), ct, reqBody)
		if err != nil {
			slog.Error(fmt.Sprintf("error posting image: %v", err))
			result.IsError = true
			result.Error = wypes.UInt32(3) // HTTPErrorRequestError
			return wypes.Void{}
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(4) // HTTPErrorRuntimeError
			return wypes.Void{}
		}

		var payload interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			result.IsError = true
			result.Error = wypes.UInt32(4) // HTTPErrorRuntimeError
			return wypes.Void{}
		}

		m := payload.(map[string]interface{})
		v, ok := m[responseKey.Unwrap()]
		if !ok {
			result.IsError = true
			result.Error = wypes.UInt32(4) // HTTPErrorRuntimeError
			return wypes.Void{}
		}

		val := v.(string)
		if len(val) >= 172 {
			val = val[:168] + "..."
		}
		res := make([]byte, len(val))
		copy(res, val)

		result.IsError = false
		result.OK = wypes.Bytes{Raw: res}
		result.DataPtr = ctx.ReturnDataPtr

		result.Lower(s)
		if s.Error != nil {
			slog.Error(fmt.Sprintf("httpPostImageFunc error in store after lower: %v", s.Error))
		}

		return wypes.Void{}
	}
}

func hostTimeNowFunc(ctx *cv.Context) func(*wypes.Store, wypes.UInt32) wypes.UInt64 {
	return func(*wypes.Store, wypes.UInt32) wypes.UInt64 {
		t := time.Now().UnixMicro()

		return wypes.UInt64(t)
	}
}
