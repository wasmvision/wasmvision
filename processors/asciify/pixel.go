//go:build tinygo

package main

import "wasmcv.org/wasm/cv/mat"

var (
	pixelLookup = []byte(" .,:;i1tfLCG08@")
	ascii       [60][80]byte
)

func imageToAscii(image mat.Mat) {
	for y := uint32(0); y < 60; y++ {
		for x := uint32(0); x < 80; x++ {
			b := image.GetUcharAt(y, x*3+0)
			g := image.GetUcharAt(y, x*3+1)
			r := image.GetUcharAt(y, x*3+2)
			a := uint8(255)

			ascii[y][x] = pixelToChar(r, g, b, a)
		}
	}
}

func pixelToChar(r, g, b, a uint8) byte {
	c := intensity(r, g, b, a)
	val := rescale(int32(c), 0, 255*3, 0, int32(len(pixelLookup)-1))

	return pixelLookup[val]
}

// full intensity for now
func intensity(r, g, b, a uint8) int32 {
	return int32(r + g + b)
}

func rescale(input, fromMin, fromMax, toMin, toMax int32) int32 {
	return (input-fromMin)*(toMax-toMin)/(fromMax-fromMin) + toMin
}
