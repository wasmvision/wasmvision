# wasmVision

wasmVision is a computer vision processing engine designed to be used with guest modules written using WebAssembly.

These processing modules can be written in Go, Rust, or the C programming language.

## How it works

wasmVision is written in the Go programming language using the [GoCV Go language wrappers](https://github.com/hybridgroup/gocv) for [OpenCV](https://github.com/opencv/opencv) and the [Wazero WASM runtime](https://github.com/tetratelabs/wazero).

wasmVision guest processing modules are compiled into WebAssembly, and expected to support the wasmCV interface. See https://github.com/hybridgroup/wasmcv

Processors can be written in Go, Rust, or the C programming language.

See the [ARCHITECTURE.md](ARCHITECTURE.md) document for more details.

## How to run it

```shell
go run ./cmd/wasmvision -module=/path/to/your/processor.wasm
```
