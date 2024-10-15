#!/bin/bash

# must have TinyGo installed to obtain the TINYGOROOT, or else separately install wasi-lib
WASI_LIBC_SYSROOT=$(tinygo env TINYGOROOT)

# must change this to match your wasmcv C files installation
WASMCV_C_COMPONENTS_PATH="../../../wasmcv/components/c"

# must change this to match your wasmcv C files installation
WASMVISION_C_COMPONENTS_PATH="../../../wasmvision-sdk/components/c"

/opt/wasi-sdk/bin/clang --target=wasm32-unknown-unknown -O3 \
        --sysroot="${WASI_LIBC_SYSROOT}/lib/wasi-libc/sysroot" \
        -z stack-size=4096 -Wl,--initial-memory=65536 \
        -I$WASMCV_C_COMPONENTS_PATH -I$WASMVISION_C_COMPONENTS_PATH \
        -o ../blurc.wasm process.c \
        $WASMCV_C_COMPONENTS_PATH/wasmcv/imports.c $WASMCV_C_COMPONENTS_PATH/wasmcv/imports_component_type.o \
        $WASMVISION_C_COMPONENTS_PATH/wasmvision/platform.c $WASMVISION_C_COMPONENTS_PATH/wasmvision/platform_component_type.o \
        -Wl,--export=process \
        -Wl,--export=__data_end -Wl,--export=__heap_base \
        -Wl,--strip-all,--no-entry \
        -Wl,--unresolved-symbols=ignore-all \
        -nostdlib \
