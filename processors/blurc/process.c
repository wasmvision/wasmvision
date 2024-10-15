#include <wasmcv/imports.h>
#include <wasmvision/platform.h>

wasm_cv_mat_own_mat_t process(wasm_cv_mat_own_mat_t image) {
    wasm_cv_cv_size_t size = {.x = 25, .y = 25};
    wasm_cv_mat_own_mat_t out_mat = wasm_cv_cv_blur(image, &size);

    platform_string_t msg = {(unsigned char *)"Blurc processor called", 23};
    wasmvision_platform_logging_log(&msg);

    return out_mat;
}
