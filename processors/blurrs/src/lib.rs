#![no_std]

extern crate core;
extern crate alloc;

use alloc::string::ToString;
use wasmcv::wasm::cv;
use wasmvision::wasmvision::platform::logging;

#[no_mangle]
pub extern fn process(mat: cv::mat::Mat) -> cv::mat::Mat {
    logging::info(&["Performing blur on image with Cols: ", &mat.cols().to_string(), " Rows: ", &mat.rows().to_string()].concat());

    if mat.empty() {
        logging::warn("image was empty");
        return mat;
    }
    
    let copy = mat.clone();
    let result = cv::cv::blur(copy, cv::types::Size{x: 25, y: 25});
    match result { 
        Ok(v) => return v, 
        Err(e) => {
            logging::error(&["Error: ", &e.to_string()].concat());
            return mat;
        },
    }
}

#[no_mangle]
pub extern fn malloc(size: usize) -> *mut u8 {
    let layout = core::alloc::Layout::from_size_align(size, 1).unwrap();
    unsafe { alloc::alloc::alloc(layout) }
}

// Use `wee_alloc` as the global allocator...for now.
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;
