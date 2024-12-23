# Architecture

## Overview

```mermaid
flowchart LR
    subgraph wasmVision
        subgraph engine
            subgraph Capture
                Devices
            end
            Runtime[WASM Runtime]
            Capture--frame-->Runtime
            Capture<-->OpenCV
            Runtime<-->OpenCV
        end
        subgraph processors
            Runtime--frame-->processor1.wasm
            Runtime--frame-->processor2.wasm
            Runtime--frame-->processor3.wasm
            Runtime--frame-->processor4.wasm
        end
    end
```

The pipeline of Processor modules are called in order, one after another. The output from the first is passed into the second, and so on. Once the last processor module has finished, the frame resources are cleaned up. Then the next frame is read from the capture device and passed into the first processor module.

### Engine

The wasmVision engine. Includes some platform capabilities such as MJPEG streaming and saving to files.

### Capture

This is how wasmVision can capture or import images or video to be processed from a connected webcam, from files, or from streams.

### Devices

Specific hardware or software devices that capture images or video. Currently supported devices are a connected camera, a file input, or a stream input using GStreamer.

### Runtime

The wasmVision runtime consists of the WebAssembly runtime engine (currently Wazero) and the wasmVision host platform API functions such as logging, configuration, and making http calls to external servers.

### Processors

The image processing modules that are used by wasmVision. These are modules written using Go, Rust, or C, and compiled into WebAssembly. They use the [wasmCV interface](https://github.com/wasmvision/wasmcv) and the [wasmVision platform interface](https://github.com/wasmvision/wasmvision-sdk). See [processors directory](./processors/) for several processors you can use, or use as examples.

### OpenCV

The computer vision processing capabilities implemented using OpenCV/GoCV. Take a look at the [`cv` directory](./cv/) for the code for this integration.
