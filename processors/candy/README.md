# candy

![candy](../../images/candy-processor.png)

wasmCV guest module that renders images using fast neural style transfer.

## How to build

```shell
tinygo build -o ../candy.wasm -target=wasm-unknown .
```

## Downloading the model

Download the file `candy-9.onnx` into your `$HOME/models` directory from https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/model/candy-9.onnx

For more information see https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/README.md
