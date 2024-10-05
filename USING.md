# Using wasmVision

```
NAME:
   wasmvision - wasmVision CLI

USAGE:
   wasmvision [global options] command [command options]

VERSION:
   0.1.0-pre5

COMMANDS:
   run       Run wasmVision processors
   download  Download computer vision models and processors
   info      Show installation info
   version   Show version
   about     About wasmVision
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

These are some of the things you can do with wasmVision.

## `wasmvision run`

```shell
NAME:
   wasmvision run - Run wasmVision processors

USAGE:
   wasmvision run [command options]

OPTIONS:
   --source value, -s value                                     video capture source to use such as webcam or file (0 is the default webcam on most systems) (default: "0")
   --output value, -o value                                     output type (mjpeg, file) (default: "mjpeg")
   --destination value, -d value                                output destination (port, file path)
   --processor value, -p value [ --processor value, -p value ]  wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm
   --logging                                                    log detailed info to console (default: true) (default: true)
   --models-dir value, --models value                           directory for model loading (default to $home/models) [$WASMVISION_MODELS_DIR]
   --model-download, --download                                 automatically download known models (default: true) (default: true)
   --processors-dir value, --processors value                   directory for processor loading (default to $home/processors) [$WASMVISION_PROCESSORS_DIR]
   --processor-download                                         automatically download known processors (default: true) (default: true)
   --help, -h                                                   show help
```

### Automatically download the `blur` processor, capture from your webcam (default), process the video, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -p blur
```

### Automatically download the `mosaic` processor and the `mosaic-9.onnx` model, capture from your webcam (default), process the video, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -p mosaic
```

### Capture from your webcam (default), process the video using the `mosaic.wasm` processor located in the default processor directory, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -p mosaic.wasm
```

### Capture from your webcam (default), process the video using the `yourcustom.wasm` processor located in a non-default directory, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -p /path/to/processors/yourcustom.wasm
```

### Capture from a secondary webcam, process the video, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -s /dev/video2 -p /path/to/processors/mosaic.wasm
```

### Capture from your webcam (default), process the video using 2 processors chained together, and stream the output using MJPEG to port 8080 (default)

```shell
wasmvision run -p /path/to/processors/hello.wasm -p /path/to/processors/mosaic.wasm
```

### Capture from a file, process the video, and stream the output using MJPEG to port 6000

```shell
wasmvision run -s /path/to/video/filename.mp4 -p /path/to/processors/blur.wasm -o mjpeg -d :6000
```

### Capture from your webcam, process the video, and save the output to a file

```shell
wasmvision run -p /path/to/processors/mosaic.wasm -o file -d /path/to/video/filename.avi
```

## `wasmvision download`

```shell
NAME:
   wasmvision download - Download computer vision models and processors

USAGE:
   wasmvision download command [command options]

COMMANDS:
   model      download a known computer vision model
   processor  download a known processor
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### Download the `candy.wasm` processor used for fast neural style transfer from the wasmVision repository to the default processors directory on the local machine.

```shell
wasmvision download processor candy
```

### Download the `candy-9.onnx` model used for fast neural style transfer from the official ONNX repository to the default models directory on the local machine.

```shell
wasmvision download model candy-9
```
