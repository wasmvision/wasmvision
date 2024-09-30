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

The wasmVision engine.

### Capture

This is how wasmVision can capture or import images or video to be processed,

### Devices

Specific hardware or software devices that capture images or video,

### Runtime

The WebAssembly runtime engine, currently Wazero.

### Processors

The wasmCV image processing modules that are used by wasmVision. See [processors](./processors/) directory.

### OpenCV

The computer vision processing capabilities implemented using OpenCV/GoCV.

