# wasmvision

wasmVision is a computer vision processing engine designed to be used with processing moduleswritten using WebAssembly.

It is a Go application written using [GoCV Go language wrappers for OpenCV](https://github.com/hybridgroup/gocv) and the [Wazero WASM runtime](https://github.com/tetratelabs/wazero).

## How it works


## How to run it

```shell
go run . -processor=/path/to/your/processor.wasm
```
