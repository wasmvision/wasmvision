package engine

import (
	"math/rand/v2"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

var NetCache = make(map[wypes.UInt32]Net)

// Net is a wrapper around gocv.Net for DNN image processing.
type Net struct {
	ID  wypes.UInt32
	Net gocv.Net
}

// NewNet creates a new Net.
func NewNet() Net {
	id := rand.IntN(102400)
	return Net{
		ID: wypes.UInt32(id),
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
