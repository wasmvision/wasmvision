#![no_std]

extern crate core;
extern crate alloc;

use alloc::string::ToString;
use wasmcv::wasm::cv;
use wasmvision::wasmvision::platform::logging;

#[no_mangle]
pub extern fn process(mat: cv::mat::Mat) -> cv::mat::Mat {
    logging::log(&["Performing blur on image with Cols: ", &mat.cols().to_string(), " Rows: ", &mat.rows().to_string()].concat());

    let out = cv::cv::blur(mat, cv::types::Size{x: 25, y: 25});
    return out;
}
