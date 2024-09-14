# wasmVision

wasmVision is a high-performance computer vision processing engine designed to be customized and extended using WebAssembly.

## How it works

```mermaid
flowchart LR
    subgraph wasmVision
        Capture-->P[WASM Processor Modules]
        P-->Output
    end
```

The wasmVision engine is written in the [Go programming language](https://go.dev/) using the [GoCV Go language wrappers](https://github.com/hybridgroup/gocv) for [OpenCV](https://github.com/opencv/opencv) and the [Wazero WASM runtime](https://github.com/tetratelabs/wazero).

wasmVision processing modules are WebAssembly guest modules that support the wasmCV interface.

See https://github.com/hybridgroup/wasmcv

These processing modules can be written in Go, Rust, or the C programming language.

See the [ARCHITECTURE.md](ARCHITECTURE.md) document for more details.

## How to run it

```shell
go run ./cmd/wasmvision -module=/path/to/your/processor.wasm
```
