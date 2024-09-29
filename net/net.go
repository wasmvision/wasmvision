package net

import (
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/orsinium-labs/wypes"
	"gocv.io/x/gocv"
)

// Net is a wrapper around gocv.Net for DNN image processing.
type Net struct {
	ID    wypes.UInt32
	Net   gocv.Net
	Model string
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

// ModelFile gets the model file path name for the Net.
func (n *Net) ModelFile() string {
	return filepath.Join(DefaultModelPath(), n.Model)
}

func DefaultModelPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, "models")
}
