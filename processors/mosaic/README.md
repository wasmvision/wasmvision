# mosaic

![mosaic](../../images/mosaic-processor.png)

wasmCV guest module that renders images in a mosaic using fast neural style transfer.

## How to build

```shell
tinygo build -o ../mosaic.wasm -target=wasm-unknown .
```

## Downloading the model

Download the file `mosaic-9.onnx` into your `$HOME/models` directory from https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/model/mosaic-9.onnx

For more information see https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/README.md
