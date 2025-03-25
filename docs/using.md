# Using wasmVision

```
NAME:
   wasmvision - wasmVision CLI

USAGE:
   wasmvision [global options] [command [command options]]

VERSION:
   0.3.0

DESCRIPTION:
   wasmVision gets you up and running with computer vision.

COMMANDS:
   run       Run wasmVision processors
   download  Download computer vision models and processors
   info      Show installation info
   listall   Lists all known models and processors
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
   wasmvision run

OPTIONS:
   --file string, -f string                                         TOML or YAML file with configuration
   --source string, -s string                                       video capture source to use. webcam id, file name, stream, testpattern (0 is the default webcam on most systems) (default: "0")
   --capture string                                                 video capture type to use (auto, ffmpeg, gstreamer, webcam) (default: "auto")
   --output string, -o string                                       output type (mjpeg, file) (default: "mjpeg")
   --destination string, -d string                                  output destination (port, file path)
   --logging string                                                 logging level to use (error, warn, info, debug) (default: "warn")
   --cuda-enable                                                    enable CUDA support (if available) (default: false)
   --processor string, -p string [ --processor string, -p string ]  wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm
   --processors-dir string                                          directory for processor loading (default to $home/processors) [$WASMVISION_PROCESSORS_DIR]
   --processor-download                                             automatically download known processors (default: true)
   --config string, -c string [ --config string, -c string ]        configuration for processors. Format: -config key1=val1 -config key2=val2
   --models-dir string                                              directory for model loading (default to $home/models) [$WASMVISION_MODELS_DIR]
   --models-download, --download                                    automatically download known models (default: true)
   --mcp-server                                                     enable MCP server (default: false)
   --mcp-port string                                                port for MCP server (default: ":5001")
   --help, -h                                                       show help
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
   wasmvision download [command [command options]] 

COMMANDS:
   model      download a known computer vision model
   processor  download a known processor

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

## `wasmvision info`

```shell
NAME:
   wasmvision info - Show installation info

USAGE:
   wasmvision info

OPTIONS:
   --help, -h  show help
```

Run this to obtain information about the installed capabilities of wasmVision.

```shell
wasmVision version 0.3.0 linux/amd64
Camera backends:  GSTREAMER V4L2 FIREWIRE UEYE OBSENSOR
Stream backends:  FFMPEG GSTREAMER INTEL_MFX V4L2 CV_IMAGES CV_MJPEG
Writer backends:  FFMPEG GSTREAMER INTEL_MFX CV_IMAGES CV_MJPEG
```

## `wasmvision listall`

```shell
NAME:
   wasmvision listall - Lists all known models and processors

USAGE:
   wasmvision listall [command [command options]] 

COMMANDS:
   models      lists all known computer vision models
   processors  lists all known wasm processors

OPTIONS:
   --help, -h  show help
```

### Show a list of all models that can be downloaded, either manually or automatically.

```shell
wasmvision listall models
```

### Show a list of all processors that can be downloaded, either manually or automatically.

```shell
wasmvision listall processors
```

## Configuration files

You can use an external TOML or YAML file to configure wasmVision. Here is an TOML example:

```toml
[main]
logging = "warn"
destination = ":8080"

[processing]
pipeline = [
    "./processors/ollama.wasm",
    "./processors/mosaic.wasm",
    "./processors/captions.wasm"
]
configuration = [
    "model=llava"
]
download = false

[models]
download = true

[server]
mcp-enable = true
mcp-port = ":5001"
```

Here is the same configuration in YAML format:

```yaml
main:
  logging: "warn"
  destination: ":8080"

processing:
  pipeline:
    - "./processors/ollama.wasm"
    - "./processors/udnie.wasm"
    - "./processors/captions.wasm"
  configuration:
    - "model=llava"
  download: false

models:
  download: true

server:
   mcp-enable: true
   mcp-port: ":5001"
```
