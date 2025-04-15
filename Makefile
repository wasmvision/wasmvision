
asciify:
	cd processors/asciify; go mod tidy; tinygo build -o ../asciify.wasm -target=wasm-unknown --no-debug .

blur:
	cd processors/blur; go mod tidy; tinygo build -o ../blur.wasm -target=wasm-unknown --no-debug .

blurrs:
	cd processors/blurrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/blurrs.wasm ../

captions:
	cd processors/captions; go mod tidy; tinygo build -o ../captions.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

edge-detect:
	cd processors/edge-detect; go mod tidy; tinygo build -o ../edge-detect.wasm -target=wasm-unknown --no-debug .

faceblur:
	cd processors/faceblur; go mod tidy; tinygo build -o ../faceblur.wasm -target=wasm-unknown --no-debug .

face-counter:
	cd processors/face-counter; go mod tidy; tinygo build -o ../face-counter.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

facedetectyn:
	cd processors/facedetectyn; go mod tidy; tinygo build -o ../facedetectyn.wasm -target=wasm-unknown --no-debug .

facedetectynrs:
	cd processors/facedetectynrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/facedetectynrs.wasm ../

face-expression:
	cd processors/face-expression; go mod tidy; tinygo build -o ../face-expression.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

gaussianblur:
	cd processors/gaussianblur; go mod tidy; tinygo build -o ../gaussianblur.wasm -target=wasm-unknown --no-debug .

hello:
	cd processors/hello; go mod tidy; tinygo build -o ../hello.wasm -target=wasm-unknown --no-debug .

object-detector:
	cd processors/object-detector; go mod tidy; tinygo build -o ../object-detector.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

ollama:
	cd processors/ollama; go mod tidy; tinygo build -o ../ollama.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

style-transfer:
	cd processors/style-transfer; go mod tidy; tinygo build -o ../style-transfer.wasm -target=wasm-unknown --no-debug .

processors: asciify blur blurrs captions edge-detect faceblur face-counter facedetectyn facedetectynrs face-expression gaussianblur hello object-detector ollama style-transfer
	@echo "All processors built successfully!"
