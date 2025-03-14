# captions

![captions](../../images/captions-processor.png)

wasmVision processor that adds captions to the final output image. 

## How to build

```shell
tinygo build -o ../captions.wasm -target=wasm-unknown .
```

## How to run

```shell
wasmvision run -p ollama.wasm -p captions.wasm
```

The `ollama.wasm` processor uses the built-in datastore to save any information about the scene description. The `captions.wasm` processor looks for this information and then displays this description on the output image.
