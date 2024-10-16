#![no_std]

extern crate core;
extern crate alloc;

use alloc::vec::Vec;
use wasmcv::wasm::cv;
use wasmvision::wasmvision::platform::logging;
use spin::{Once, Mutex};

static INIT: Once = Once::new();
static mut DETECTOR: Option<Mutex<cv::objdetect::FaceDetectorYn>> = None;

const RED: cv::types::Rgba = cv::types::Rgba { r: 255, g: 0, b: 0, a: 255 };
const BLUE: cv::types::Rgba = cv::types::Rgba { r: 0, g: 0, b: 255, a: 255 };
const GREEN: cv::types::Rgba = cv::types::Rgba { r: 0, g: 255, b: 0, a: 255 };
const PINK: cv::types::Rgba = cv::types::Rgba { r: 255, g: 0, b: 255, a: 255 };
const YELLOW: cv::types::Rgba = cv::types::Rgba { r: 255, g: 255, b: 0, a: 255 };

fn init() {
    INIT.call_once(|| {
        unsafe {
            DETECTOR = Some(Mutex::new(cv::objdetect::FaceDetectorYn::new(
                "face_detection_yunet_2023mar", "", cv::types::Size { x: 200, y: 200 }
            )));
        }
    });
}

#[no_mangle]
pub extern "C" fn _start() {
    init();
}

#[no_mangle]
pub extern fn process(mat: cv::mat::Mat) -> cv::mat::Mat {
    if mat.empty() {
        return mat;
    }

    let mut out = mat.clone();
    let sz = out.size();
    let sz_i32: Vec<i32> = sz.into_iter().map(|x| x as i32).collect();

    unsafe {
        if let Some(ref detector) = DETECTOR {
            let detector = detector.lock();
            detector.set_input_size(cv::types::Size { x: sz_i32[1], y: sz_i32[0] });

            let faces = detector.detect(mat.clone());
            draw_faces(&mut out, &faces);
        }
    }

    logging::log("Performed face detection on image");

    out
}

fn draw_faces(out: &mut cv::mat::Mat, faces: &cv::mat::Mat) {
    for r in 0..faces.rows() {
        let x0 = faces.get_float_at(r, 0) as i32;
        let y0 = faces.get_float_at(r, 1) as i32;
        let x1 = x0 + (faces.get_float_at(r, 2) as i32);
        let y1 = y0 + (faces.get_float_at(r, 3) as i32);

        let face_rect = cv::types::Rect { min: cv::types::Size { x: x0, y: y0 }, max: cv::types::Size { x: x1, y: y1 } };

        let right_eye = cv::types::Size { x: faces.get_float_at(r, 4) as i32, y: faces.get_float_at(r, 5) as i32 };
        let left_eye = cv::types::Size { x: faces.get_float_at(r, 6) as i32, y: faces.get_float_at(r, 7) as i32 };
        let nose_tip = cv::types::Size { x: faces.get_float_at(r, 8) as i32, y: faces.get_float_at(r, 9) as i32 };
        let right_mouth_corner = cv::types::Size { x: faces.get_float_at(r, 10) as i32, y: faces.get_float_at(r, 11) as i32 };
        let left_mouth_corner = cv::types::Size { x: faces.get_float_at(r, 12) as i32, y: faces.get_float_at(r, 13) as i32 };

        cv::cv::rectangle(out, face_rect, GREEN, 1);
        cv::cv::circle(out, right_eye, 1, BLUE, 1);
        cv::cv::circle(out, left_eye, 1, RED, 1);
        cv::cv::circle(out, nose_tip, 1, GREEN, 1);
        cv::cv::circle(out, right_mouth_corner, 1, PINK, 1);
        cv::cv::circle(out, left_mouth_corner, 1, YELLOW, 1);
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
