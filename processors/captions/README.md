# captions

wasmVision processor that adds text captions to the final output image.

## How to build

```shell
tinygo build -o ../captions.wasm -target=wasm-unknown --no-debug .
```

## How to run

```shell
wasmvision run -p ollama.wasm -p captions.wasm
```

The `ollama.wasm` processor uses the built-in datastore to save any information about the scene description. The `captions.wasm` processor looks for this information and then displays this description on the output image.

## Configuration

The following configuration settings are available for the `captions.wasm` processor.

### `caption-words-per-line`

Set the number of words per line for the captioning like this:

```shell
-c caption-words-per-line=7
```

Default value: 10

### `caption-line-height`

Set the line height for the captioning like this:

```shell
-c caption-line-height=30
```

Default value: 20

### `caption-size`

Set the font size for the captioning like this:

```shell
-c caption-size=0.8
```

Default value: 0.6

### `caption-color`

Set the font color for the captioning like this:

```shell
-c caption-color=red
```

Default value: "black"
