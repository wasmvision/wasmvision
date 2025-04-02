package cv

import (
	"log/slog"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

// Net is a wrapper around gocv.Net for DNN image processing.
type Net struct {
	ID       wypes.UInt32
	Name     string
	Filename string

	Net gocv.Net
}

// NewNet creates a new Net.
func NewNet(model string) *Net {
	return &Net{
		Name: model,
	}
}

// SetNet sets the gocv.Net for the Net.
func (n *Net) SetNet(net gocv.Net) {
	n.Net = net
}

// Close closes the Net.
func (n *Net) Close() {
	n.Net.Close()
}

// Layer is a wrapper around gocv.Layer for DNN image processing.
type Layer struct {
	ID wypes.UInt32

	Layer gocv.Layer
}

func NewLayer(layer gocv.Layer) *Layer {
	return &Layer{
		Layer: layer,
	}
}

func (l *Layer) SetLayer(layer gocv.Layer) {
	l.Layer = layer
}

func (l *Layer) Close() {
	l.Layer.Close()
}

func handleNetReturn(ctx *Context, s *wypes.Store, net *Net, result wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32]) {
	result.IsError = false
	result.OK = wypes.HostRef[*Net]{Raw: net}
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}

func handleNetError(ctx *Context, s *wypes.Store, net *Net, result wypes.Result[wypes.HostRef[*Net], wypes.HostRef[*Net], wypes.UInt32], err error) {
	if err == nil {
		return
	}
	if net != nil {
		net.Close()
	}

	slog.Error("cv net error", "error", err)
	s.Error = err
	result.IsError = true
	result.Error = 1
	result.DataPtr = ctx.ReturnDataPtr
	result.Lower(s)
}
