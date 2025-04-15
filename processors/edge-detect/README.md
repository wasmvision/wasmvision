# edge-detect

![edge-detect](../../images/edge-detect-processor.png)

wasmVision processor that detects edges using a model called Dense Extreme Inception Network for Edge Detection (DexiNed), a Convolutional Neural Network (CNN) architecture for edge detection.

## How to build

```shell
tinygo build -o ../edge-detect.wasm -target=wasm-unknown .
```

## Downloading the model

The first time you run the processor it will automatically download the model, or you can download it by running the command:

```shell
wasmvision download dexined_2024sep
```

For more information see:

https://github.com/opencv/opencv_zoo/blob/main/models/edge_detection_dexined/README.md 

and also 

https://github.com/xavysp/DexiNed
