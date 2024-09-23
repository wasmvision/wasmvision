# Development

It you want to try out the latest development builds, or work on developing wasmVision itself, this is the place for information.

## Latest builds

### Linux

You can download the latest development builds for Linux by looking under [Linux workflows](https://github.com/wasmvision/wasmvision/actions/workflows/linux.yml) then clicking on the latest one you want to test. 

Under the "Artifacts" click on the link for either "wasmvision-linux-amd64" or "wasmvision-linux-arm64" depending on your processor to download it.

Extract the file to the desired test directory.

### macOS

You can download the latest development builds for macOs by looking under [macOS workflows](https://github.com/wasmvision/wasmvision/actions/workflows/macos.yml) then clicking on the latest one you want to test. 

Under the "Artifacts" click on the link for "wasmvision-macos-arm64" to download it.

Extract the file to the desired test directory.

### Windows

You can download the latest development builds for Linux by looking under [Linux workflows](https://github.com/wasmvision/wasmvision/actions/workflows/linux.yml) then clicking on the latest one you want to test. 

Under the "Artifacts" click on the link for "wasmvision-windows-amd64" to download it.

NOTE: you will likely need to configure your Windows Defender to download the ZIP file with the `wasmvision.exe` executable.

Extract the executable to your desired directory.

### Docker

You can run the latest development builds of wasmVision using Docker.

Pull the current development version:

```shell
docker pull ghcr.io/wasmvision/wasmvision:main
```

Run your desired docker commands using the tagged image `ghcr.io/wasmvision/wasmvision:main` (NOT `ghcr.io/wasmvision/wasmvision:latest` which is the latest released version).

## Local development

If you have a local installation of both Go and OpenCV you can install wasmVision directly:

```shell
git clone https://github.com/wasmvision/wasmvision.git
cd wasmvision
go install ./cmd/wasmvision/
```

And run it:

```shell
wasmvision run -p ./processors/hello.wasm -mjpeg=true


You can run wasmVision using Docker.

Pull the current development version:

```shell
docker pull ghcr.io/wasmvision/wasmvision:main
```
