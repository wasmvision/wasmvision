
asciify:
	cd processors/asciify; tinygo build -o ../asciify.wasm -target=wasm-unknown .

blur:
	cd processors/blur; tinygo build -o ../blur.wasm -target=wasm-unknown .

blurrs:
	cd processors/blurrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/blurrs.wasm ../

candy:
	cd processors/candy; tinygo build -o ../candy.wasm -target=wasm-unknown .

faceblur:
	cd processors/faceblur; tinygo build -o ../faceblur.wasm -target=wasm-unknown .

facedetectyn:
	cd processors/facedetectyn; tinygo build -o ../facedetectyn.wasm -target=wasm-unknown .

gaussianblur:
	cd processors/gaussianblur; tinygo build -o ../gaussianblur.wasm -target=wasm-unknown .

hello:
	cd processors/hello; tinygo build -o ../hello.wasm -target=wasm-unknown .

ollama:
	cd processors/ollama; tinygo build -o ../ollama.wasm -target=wasm-unknown .

mosaic:
	cd processors/mosaic; tinygo build -o ../mosaic.wasm -target=wasm-unknown .

pointilism:
	cd processors/pointilism; tinygo build -o ../pointilism.wasm -target=wasm-unknown .

rain-princess:
	cd processors/rainprincess; tinygo build -o ../rainprincess.wasm -target=wasm-unknown .

udnie:
	cd processors/udnie; tinygo build -o ../udnie.wasm -target=wasm-unknown .

processors: asciify blur blurrs candy faceblur facedetectyn gaussianblur hello ollama mosaic pointilism rain-princess udnie
