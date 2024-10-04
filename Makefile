
asciify:
	cd processors/asciify; tinygo build -o ../asciify.wasm -target=wasm-unknown .

blur:
	cd processors/blur; tinygo build -o ../blur.wasm -target=wasm-unknown .

blurrs:
	cd processors/blurrs; cargo build --target wasm32-unknown-unknown --release; \
		cp ./target/wasm32-unknown-unknown/release/blurrs.wasm ../

candy:
	cd processors/candy; tinygo build -o ../candy.wasm -target=wasm-unknown .

gaussianblur:
	cd processors/gaussianblur; tinygo build -o ../gaussianblur.wasm -target=wasm-unknown .

hello:
	cd processors/hello; tinygo build -o ../hello.wasm -target=wasm-unknown .

mosaic:
	cd processors/mosaic; tinygo build -o ../mosaic.wasm -target=wasm-unknown .

processors: asciify blur blurrs candy gaussianblur hello mosaic
