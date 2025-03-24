
asciify:
	cd processors/asciify; go mod tidy; tinygo build -o ../asciify.wasm -target=wasm-unknown --no-debug .

blur:
	cd processors/blur; go mod tidy; tinygo build -o ../blur.wasm -target=wasm-unknown  --no-debug .

blurrs:
	cd processors/blurrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/blurrs.wasm ../

candy:
	cd processors/candy; go mod tidy; tinygo build -o ../candy.wasm -target=wasm-unknown --no-debug .

captions:
	cd processors/captions; go mod tidy; tinygo build -o ../captions.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

faceblur:
	cd processors/faceblur; go mod tidy; tinygo build -o ../faceblur.wasm -target=wasm-unknown --no-debug .

facedetectyn:
	cd processors/facedetectyn; go mod tidy; tinygo build -o ../facedetectyn.wasm -target=wasm-unknown --no-debug .

facedetectynrs:
	cd processors/facedetectynrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/facedetectynrs.wasm ../

gaussianblur:
	cd processors/gaussianblur; go mod tidy; tinygo build -o ../gaussianblur.wasm -target=wasm-unknown --no-debug .

hello:
	cd processors/hello; go mod tidy; tinygo build -o ../hello.wasm -target=wasm-unknown --no-debug .

ollama:
	cd processors/ollama; go mod tidy; tinygo build -o ../ollama.wasm -target=wasip1 -buildmode=c-shared -scheduler=none --no-debug .

mosaic:
	cd processors/mosaic; go mod tidy; tinygo build -o ../mosaic.wasm -target=wasm-unknown --no-debug .

pointilism:
	cd processors/pointilism; go mod tidy; tinygo build -o ../pointilism.wasm -target=wasm-unknown --no-debug .

rain-princess:
	cd processors/rainprincess; go mod tidy; tinygo build -o ../rainprincess.wasm -target=wasm-unknown --no-debug .

udnie:
	cd processors/udnie; go mod tidy; tinygo build -o ../udnie.wasm -target=wasm-unknown --no-debug .

processors: asciify blur blurrs candy captions faceblur facedetectyn facedetectynrs gaussianblur hello ollama mosaic pointilism rain-princess udnie
