# style-transfer

![udnie](../../images/udnie-processor.png)

wasmVision processor that renders images using fast neural style transfer models.

If GPU is available and enabled, this processor will automatically use hardware acceleration.

## How to build

```shell
tinygo build -o ../style-transfer.wasm -target=wasm-unknown --no-debug .
```

## How to run

```shell
wasmvision run -p style-transfer.wasm -c style-model=udnie
```

The `style-transfer.wasm` processor uses the selected model to transform the input into the output image.

## Configuration

The following configuration settings are available for the `captions.wasm` processor.

### `style-model`

Set the fast neural style trasfer model like this:

```shell
-c style-model="candy-9"
```

Default value: "mosaic"

## Fast Neural Style Transfer models

The following models work with this processor:

### `candy`

![candy](../../images/candy-processor.png)

```shell
wasmvision run -p style-transfer.wasm -c style-model="candy-9"
```

### `mosaic`

![mosaic](../../images/mosaic-processor.png)

```shell
wasmvision run -p style-transfer.wasm -c style-model="mosaic-9"
```

### `pointilism`

![pointilism](../../images/pointilism-processor.png)

```shell
wasmvision run -p style-transfer.wasm -c style-model="pointilism-9"
```

### `rainprincess`

![rain princess](../../images/rainprincess-processor.png)

```shell
wasmvision run -p style-transfer.wasm -c style-model="rain-princess-9"
```

### `udnie`

![udie](../../images/udnie-processor.png)

```shell
wasmvision run -p style-transfer.wasm -c style-model="udnie-9"
```

## Downloading models

The first time you run the processor it will automatically download the selected model, or you can download it by running the command:

```shell
wasmvision download candy-9
```

To see a list of available models:

```shell
wasmvision listall models
```

For more information see https://github.com/onnx/models/blob/main/validated/vision/style_transfer/fast_neural_style/README.md
