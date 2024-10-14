# ollama

![ollama](../../images/ollama-processor.png)

wasmVision processor that sends image frames to an [Ollama](https://ollama.com/) server running a model for image description such as `llava`.

LLaVA is a multimodal model that combines a vision encoder and Vicuna for general-purpose visual and language understanding, achieving impressive chat capabilities mimicking spirits of the multimodal GPT-4.

## How to build

```shell
tinygo build -o ../ollama.wasm -target=wasm-unknown .
```

## Running Ollama

In order to run this processor, you need to first run the Ollama server with the `llava` model.

For example:

```shell
docker run --gpus=all -d -v ${HOME}/.ollama:/root/.ollama -v ${HOME}/ollama-import:/root/ollama-import -p 11434:11434 --name ollama ollama/ollama:latest
docker exec ollama ollama llava
```

For more information see https://ollama.com/library/llava:13b
