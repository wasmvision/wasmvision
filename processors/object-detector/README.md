# object-detector

![object-detector](../../images/object-detector-processor.png)

wasmVision processor that performs object detection using the YOLOv8 real-time object detection model. YOLOv8 was released by Ultralytics on January 10th, 2023, offering cutting-edge performance in terms of accuracy and speed.

The YOLOv8 series offers a diverse range of models, each specialized for specific tasks in computer vision. These models are designed to cater to various requirements, from object detection to more complex tasks like instance segmentation, pose/keypoints detection, oriented object detection, and classification.

wasmVision currently only provides support for the YOLOv8 detection model.

## How to build

```shell
tinygo build -o ../object-detector.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .
```

## Configuration

The following configuration settings are available for the `object-detector.wasm` processor.

### `yolo-model`

Set which YOLOv8 model to use like this:

```shell
-c yolo-model="yolov8m"
```

Default value: "yolov8n"

## Downloading models

The first time you run the processor it will automatically download the model, or you can download it by running the command:

```shell
wasmvision download yolov8n
```

To see a list of available models:

```shell
wasmvision listall models
```

For more information see:

https://docs.ultralytics.com/models/yolov8/
