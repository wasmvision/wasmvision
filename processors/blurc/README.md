# blur

![blur](../../images/blur-processor.png)

wasmVision processor in C that performs a blur.

## How to build

### wasi-sdk

Download and install `wasi-sdk` from https://github.com/WebAssembly/wasi-sdk/releases


### wasmcv

Download `wasmcv` C components:

```shell
cd <your_sources_directory>
git clone https://github.com/wasmvision/wasmcv.git 
```

### wasmvsion C processor

set `WASMCV_C_COMPONENTS_PATH` variable in `<wasmvision_directory>/processors/blurc/build.sh`

```shell
WASMCV_C_COMPONENTS_PATH="<your_sources_directory>/wasmcv/components/c"
```


### wasmVision C SDK

Download `wasmvision-sdk` C components:

```shell
cd <your_sources_directory>
git clone https://github.com/wasmvision/wasmvision-sdk.git 
```

### wasmVision C processor

set `WASMVISION_C_COMPONENTS_PATH` variable in `<wasmvision_directory>/processors/blurc/build.sh`

```shell
WASMVISION_C_COMPONENTS_PATH="<your_sources_directory>/wasmvision-sdk/components/c"
```

### Building

run:

```shell
./build.sh
```

Note that you need to also have [TinyGo](https://github.com/tinygo-org/tinygo/releases) installed to obtain the needed `wasi-libc` or else install it separately.
