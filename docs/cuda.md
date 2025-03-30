# Using wasmVision with CUDA

wasmVision includes experimental support for GPU acceleration using [CUDA](https://en.wikipedia.org/wiki/CUDA).

The initial implementation is only to [accelerate processing of DNN models](https://github.com/opencv/opencv/wiki/Deep-Learning-in-OpenCV).

In addition, initial support is Linux amd64 support only.

## Local installation

### Prerequisites

To build/run locally you must install OpenCV with CUDA support.

For more information see:

https://github.com/hybridgroup/gocv/blob/release/cuda/README.md#installing-cuda

### Building

Once you have installed the needed versions of CUDA and OpenCV, actually building the executable is much easier.

```shell
go install -tags cuda ./cmd/wasmvision
```

### Running

```shell
wasmvision run -p mosaic --cuda-enable=true
```

## Docker

You can also run a wasmVision Docker image with CUDA support.

### Prerequisites

You must install the following software before running the wasmVision Docker image with CUDA support:

- [CUDA version 12](https://docs.nvidia.com/cuda/archive/12.6.0/cuda-installation-guide-linux/index.html)

- [Nvidia Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html) 

### Running

Pull the current development version:

```shell
docker pull ghcr.io/wasmvision/wasmvision-cuda-12:main
```

Verify it is installed like this:

```shell
docker run ghcr.io/wasmvision/wasmvision-cuda-12:main version
```

Now you can run a test to capture video using your webcam, blur it using a WebAssembly processor, and then stream the output to port 8080 on your local machine:

```shell
docker run --privileged --network=host --gpus all ghcr.io/wasmvision/wasmvision-cuda-12:main run -p /processors/mosaic.wasm --cuda-enable=true
```

Point your browser to `http://localhost:8080` and you can see the output.
