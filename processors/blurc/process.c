#include <wasmcv/imports.h>
#include <wasmvision/platform.h>

wasm_cv_mat_own_mat_t process(wasm_cv_mat_own_mat_t image) {
    wasm_cv_cv_size_t size = {.x = 25, .y = 25};
    wasm_cv_cv_own_mat_t out_mat;
    wasm_cv_cv_error_result_t err;

    if (!wasm_cv_cv_blur(image, &size, &out_mat, &err)) {
        platform_string_t msg = {(unsigned char *)"Blurc processor error", 21};
        wasmvision_platform_logging_error(&msg);
        return image;
    }

    platform_string_t msg = {(unsigned char *)"Blurc processor called", 23};
    wasmvision_platform_logging_debug(&msg);

    return out_mat;
}

extern unsigned char __heap_base;

unsigned char* bump_pointer = &__heap_base;
void* malloc(int n) {
  unsigned char* r = bump_pointer;
  bump_pointer += n;
  return (void *)r;
}
