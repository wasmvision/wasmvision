# facedetectynrs

![facedetectyn](../../images/facedetectyn-processor.png)

wasmVision processor written using Rust that recognizes faces using YuNet, a light-weight, fast and accurate face detection model, which achieves 0.834(AP_easy), 0.824(AP_medium), 0.708(AP_hard) on the WIDER Face validation set.

## How to build

```shell
cargo build --target wasm32-unknown-unknown --release
cp ./target/wasm32-unknown-unknown/release/facedetectynrs.wasm ../
```
