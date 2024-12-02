unzip wasmvision-linux-arm64.zip
chmod +x wasmvision
tar -czvf wasmvision-linux-arm64.tar.gz wasmvision
rm wasmvision
rm wasmvision-linux-arm64.zip

unzip wasmvision-linux-amd64.zip
chmod +x wasmvision
tar -czvf wasmvision-linux-amd64.tar.gz wasmvision
rm wasmvision
rm wasmvision-linux-amd64.zip

unzip wasmvision-macos-arm64.zip
chmod +x wasmvision
tar -czvf wasmvision-macos-arm64.tar.gz wasmvision
rm wasmvision
rm wasmvision-macos-arm64.zip
sha256sum wasmvision-macos-arm64.tar.gz
