# Architecture

## Overview

```mermaid
flowchart LR
    subgraph wasmVision
        Capture--frame-->Runtime[WASM Runtime]
        Capture<-->Devices
        Devices<-->OpenCV
        Runtime<-->wasmCV
        Runtime<-->OpenCV
    end
    subgraph wasmCV
        processor1.wasm--frame-->processor2.wasm
        processor2.wasm--frame-->processor3.wasm
        processor3.wasm--frame-->processor4.wasm
    end
```

### Engine

The host application.

### Capture

This is how wasmVision can capture or import images or video to be processed,

### Devices

Specific hardware or software devices that capture images or video,

### Runtime

The WebAssembly runtime engine, currently Wazero.

### Modules

The wasmCV image processing modules that developers are writing.

### CV

The computer vision processing capabilities implemented using OpenCV/GoCV.

