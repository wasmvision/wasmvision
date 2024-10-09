# faceblur

![faceblur](../../images/faceblur-processor.png)

wasmVision processor that blurs faces by first detecting them and then blurring each before outputting the final image. 

The detection is done using YuNet, a light-weight, fast and accurate face detection model, which achieves 0.834(AP_easy), 0.824(AP_medium), 0.708(AP_hard) on the WIDER Face validation set.

## How to build

```shell
tinygo build -o ../faceblur.wasm -target=wasm-unknown .
```

## Downloading the model

The first time you run the processor it will automatically download the model, or you can download it by running the command:

```shell
wasmvision download face_detection_yunet_2023mar
```

For more information see https://github.com/opencv/opencv_zoo/blob/main/models/face_detection_yunet/README.md
