#include <string.h>
#include <wasmcv/imports.h>


wasm_cv_mat_borrow_mat_t process(wasm_cv_mat_borrow_mat_t image) {
    wasm_cv_cv_own_mat_t frame;
    wasm_cv_cv_own_mat_t out_mat;
    wasm_cv_cv_size_t size = {.x = 25, .y = 25};
       
    frame = wasm_cv_mat_method_mat_clone(image);
    wasm_cv_mat_method_mat_close(image);

    out_mat = wasm_cv_cv_blur(frame, &size);

    image = wasm_cv_mat_borrow_mat(out_mat);

    return image;
}
