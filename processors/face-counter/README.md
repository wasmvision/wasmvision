# face-counter

wasmVision processor that counts the average number of previously detected faces, and saves it to the Processor datastore. 

The detection is first done using the `facedetectyn.wasm` processor. Then, this processor counts any faces found for each frame and saves the average number of faces seen in JSON format:

```json
{"timestamp":"2025-04-15T17:53:08.325458Z","average-faces-seen":1}
```

The frequency of the data recorded can be changed using a configuration setting.

## How to build

```shell
tinygo build -o ../face-counter.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .
```

## How to run

```shell
wasmvision run -p facedetectyn.wasm -p face-counter.wasm
```

The `facedetectyn.wasm` processor uses the built-in datastore to save any information about detected faces for each frame. The `face-counter.wasm` processor looks for this information and then uses it to determine the average number of faces seen per minute.

## Configuration

The following configuration settings are available for the `face-counter.wasm` processor.

### `face-counter-frequency`

Set the frequency in seconds to record the data like this:

```shell
-c face-counter-frequency=60
```

Default value: 10
