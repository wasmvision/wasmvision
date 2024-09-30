# candy

![candy](../../images/candy-processor.png)

wasmVision processor that renders images using fast neural style transfer.

## How to build

```shell
tinygo build -o ../candy.wasm -target=wasm-unknown .
```

## Downloading the model

```shell
wasmvision download candy-9
```

For more information see https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/README.md
