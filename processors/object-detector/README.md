# object-detector

![object-detector](../../images/object-detector-processor.png)

wasmVision processor that performs object detection using the YOLOv8 real-time object detection model. YOLOv8 was released by Ultralytics on January 10th, 2023, offering cutting-edge performance in terms of accuracy and speed.

The YOLOv8 series offers a diverse range of models, each specialized for specific tasks in computer vision. These models are designed to cater to various requirements, from object detection to more complex tasks like instance segmentation, pose/keypoints detection, oriented object detection, and classification.

wasmVision currently only provides support for the YOLOv8 detection model.

## How to build

```shell
tinygo build -o ../object-detector.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .
```

## Downloading the model

The first time you run the processor it will automatically download the model, or you can download it by running the command:

```shell
wasmvision download yolov8n
```

For more information see:

https://docs.ultralytics.com/models/yolov8/
