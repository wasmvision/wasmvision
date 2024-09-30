# mosaic

![mosaic](../../images/mosaic-processor.png)

wasmVision processor that renders images in a mosaic using fast neural style transfer.

## How to build

```shell
tinygo build -o ../mosaic.wasm -target=wasm-unknown .
```

## Downloading the model

```shell
wasmvision download mosaic-9
```

For more information see https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/README.md
