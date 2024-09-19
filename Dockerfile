# wasmvision Dockerfile
#
# This Dockerfile is used to build a Docker image containing the wasmvision
# binary. The resulting image can be run using the `docker run` command.
#
# Example:
#
#   docker run --privileged --network=host --platform=linux/arm64  wasmvision:dev -processors=/examples/processors/blur.wasm -mjpeg=true
#
# to build this docker image:
#   docker buildx build -t wasmvision:dev --platform=linux/amd64,linux/arm64 .
#
# first stage: build the wasmvision binary
FROM --platform=$BUILDPLATFORM ghcr.io/hybridgroup/opencv:4.10.0-static AS wasmvision-build

COPY . /src

WORKDIR /src

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -tags static -o /build/wasmvision ./cmd/wasmvision

# second stage: create a minimal image with the wasmvision binary
FROM ubuntu:22.04 AS wasmvision-final

COPY --from=wasmvision-build /build/wasmvision /run/wasmvision

COPY ./processors/*.wasm /processors/

ENTRYPOINT ["/run/wasmvision"]
