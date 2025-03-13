# faceblur

![faceblur](../../images/faceblur-processor.png)

wasmVision processor that blurs previously detected faces before outputting the final image. 

The detection is first done using the `facedetectyn` processor. Then this processor acts on any faces found for each frame and blurs them for the output frame.

## How to build

```shell
tinygo build -o ../faceblur.wasm -target=wasm-unknown .
```

## How to run

```shell
wasmvision run -p facedetectyn.wasm -p faceblur.wasm
```

The `facedetectyn.wasm` processor uses the built-in datastore to save any information about detected faces for each frame. The `faceblur.wasm` processor looks for this information and then blurs any faces that have been detected for each frame for the final output image.
