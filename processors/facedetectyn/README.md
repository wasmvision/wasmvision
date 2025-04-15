# facedetectyn

![facedetectyn](../../images/facedetectyn-processor.png)

wasmVision processor that recognizes faces using YuNet, a light-weight, fast and accurate face detection model, which achieves 0.834(AP_easy), 0.824(AP_medium), 0.708(AP_hard) on the WIDER Face validation set.

## How to build

```shell
tinygo build -o ../facedetectyn.wasm -target=wasm-unknown .
```

## Downloading the model

The first time you run the processor it will automatically download the model, or you can download it by running the command:

```shell
wasmvision download yunet_2023mar
```

For more information see https://github.com/opencv/opencv_zoo/blob/main/models/face_detection_yunet/README.md

## Configuration

The following configuration settings are available for the `facedetectyn.wasm` processor.

### `detect-draw-faces`

Turn on/off the drawing of face rects like this:

```shell
-c detect-draw-faces=false
```

Default value: "true"

The data for detected faces will still be saved in the frame datastore. This setting is to make it possible to perform further processing in downstream processors based on the unmodified image data.
