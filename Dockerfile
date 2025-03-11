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
FROM --platform=${TARGETPLATFORM} ghcr.io/hybridgroup/opencv:4.11.0-alpine-ffmpeg-gstreamer AS opencv

RUN apk update && apk add --no-cache \
    util-linux-static util-linux-dev build-base \
    cmake \
    git \
    wget \
    unzip \
    pkgconfig \
    glib-static glib-dev \
    gobject-introspection-dev libmount libeconf-dev

# Install Go
ARG GO_VERSION=1.24.1
ARG TARGETARCH

RUN wget https://golang.org/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${TARGETARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${TARGETARCH}.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go

COPY . /src
WORKDIR /src

# second stage: build the wasmvision binary (linux/amd64)
FROM --platform=linux/amd64 opencv AS wasmvision-build-amd64

ENV CGO_CXXFLAGS="--std=c++11"
ENV CGO_CPPFLAGS="-I/usr/local/include/opencv4"
ENV CGO_LDFLAGS="-static -L/usr/local/lib -lopencv_gapi -lopencv_stitching -lopencv_alphamat -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dnn_superres -lopencv_dpm -lopencv_face -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_intensity_transform -lopencv_line_descriptor -lopencv_mcc -lopencv_quality -lopencv_rapid -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_signal -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_superres -lopencv_optflow -lopencv_surface_matching -lopencv_tracking -lopencv_highgui -lopencv_datasets -lopencv_text -lopencv_plot -lopencv_videostab -lopencv_videoio -lopencv_wechat_qrcode -lopencv_xfeatures2d -lopencv_shape -lopencv_ml -lopencv_ximgproc -lopencv_video -lopencv_xobjdetect -lopencv_objdetect -lopencv_calib3d -lopencv_imgcodecs -lopencv_features2d -lopencv_dnn -lopencv_flann -lopencv_xphoto -lopencv_photo -lopencv_imgproc -lopencv_core -L/usr/local/lib/opencv4/3rdparty -llibprotobuf -lade -littnotify -llibwebp -llibtiff -llibopenjp2 -lippiw -lippicv -llibjpeg-turbo -llibpng -L/lib -lzlib -lIlmImf -ldl -lm -lpthread -lrt -lavdevice -lm -latomic -lavfilter -pthread -lm -latomic -lswscale -lm -latomic -lpostproc -lm -latomic -lavformat -lm -latomic -lavcodec -lvpx -lx264 -lswresample -lm -latomic -lavutil -lbz2_static -llzma -lgstreamer-full-1.0 -lgstadaptivedemux-1.0 -lgstallocators-1.0 -lgstanalytics-1.0 -lgstapp-1.0 -lgstaudio-1.0 -lgstbadaudio-1.0 -lgstbase-1.0 -lgstbasecamerabinsrc-1.0 -lgstcodecparsers-1.0 -lgstcodecs-1.0 -lgstcontroller-1.0 -lgstcuda-1.0 -lgstfft-1.0 -lgstinsertbin-1.0 -lgstisoff-1.0 -lgstmpegts-1.0 -lgstmse-1.0 -lgstnet-1.0 -lgstpbutils-1.0 -lgstphotography-1.0 -lgstplay-1.0 -lgstplayer-1.0  -lgstreamer-1.0 -lgstriff-1.0 -lgstrtp-1.0 -lgstrtsp-1.0 -lgstsctp-1.0 -lgstsdp-1.0 -lgsttag-1.0 -lgsttranscoder-1.0 -lgsturidownloader-1.0 -lgstvideo-1.0 -lgstwebrtc-1.0 -lopenh264 -lstdc++ -lglib-2.0 -lgio-2.0 -lgmodule-2.0 -lgobject-2.0 -lgthread-2.0 -lgirepository-2.0 -lffi -lpcre2-8 -lintl -L/usr/local/lib/gstreamer-1.0 -lgstcoreelements -lgstapp -lgstplayback -lgstrawparse -lgsttcp -lgstvideoconvertscale -lgstvideotestsrc -lgstaudioparsers -lgstcodectimestamper  -lgstisomp4 -lgstopenh264 -lgstrtp -lgstrtpmanager -lgstrtsp -lgstudp -lgstvideoparsersbad -lmount -lblkid -leconf"


# second stage: build the wasmvision binary (linux/amd64)
FROM --platform=linux/arm64 opencv AS wasmvision-build-arm64

ENV CGO_CXXFLAGS="--std=c++11"
ENV CGO_CPPFLAGS="-I/usr/local/include/opencv4"
ENV CGO_LDFLAGS="-static -L/usr/local/lib -lopencv_gapi -lopencv_stitching -lopencv_alphamat -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dnn_superres -lopencv_dpm -lopencv_face -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_intensity_transform -lopencv_line_descriptor -lopencv_mcc -lopencv_quality -lopencv_rapid -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_signal -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_superres -lopencv_optflow -lopencv_surface_matching -lopencv_tracking -lopencv_highgui -lopencv_datasets -lopencv_text -lopencv_plot -lopencv_videostab -lopencv_videoio -lopencv_wechat_qrcode -lopencv_xfeatures2d -lopencv_shape -lopencv_ml -lopencv_ximgproc -lopencv_video -lopencv_xobjdetect -lopencv_objdetect -lopencv_calib3d -lopencv_imgcodecs -lopencv_features2d -lopencv_dnn -lopencv_flann -lopencv_xphoto -lopencv_photo -lopencv_imgproc -lopencv_core -L/usr/local/lib/opencv4/3rdparty -llibprotobuf -lade -littnotify -llibwebp -llibtiff -llibopenjp2 -llibjpeg-turbo -llibpng -L/lib -lzlib -lIlmImf -ldl -lm -lpthread -lrt -lavdevice -lm -latomic -lavfilter -pthread -lm -latomic -lswscale -lm -latomic -lpostproc -lm -latomic -lavformat -lm -latomic -lavcodec -lvpx -lx264 -lswresample -lm -latomic -lavutil -lbz2_static -llzma -lgstreamer-full-1.0 -lgstadaptivedemux-1.0 -lgstallocators-1.0 -lgstanalytics-1.0 -lgstapp-1.0 -lgstaudio-1.0 -lgstbadaudio-1.0 -lgstbase-1.0 -lgstbasecamerabinsrc-1.0 -lgstcodecparsers-1.0 -lgstcodecs-1.0 -lgstcontroller-1.0 -lgstcuda-1.0 -lgstfft-1.0 -lgstinsertbin-1.0 -lgstisoff-1.0 -lgstmpegts-1.0 -lgstmse-1.0 -lgstnet-1.0 -lgstpbutils-1.0 -lgstphotography-1.0 -lgstplay-1.0 -lgstplayer-1.0  -lgstreamer-1.0 -lgstriff-1.0 -lgstrtp-1.0 -lgstrtsp-1.0 -lgstsctp-1.0 -lgstsdp-1.0 -lgsttag-1.0 -lgsttranscoder-1.0 -lgsturidownloader-1.0 -lgstvideo-1.0 -lgstwebrtc-1.0 -lopenh264 -lstdc++ -lglib-2.0 -lgio-2.0 -lgmodule-2.0 -lgobject-2.0 -lgthread-2.0 -lgirepository-2.0 -lffi -lpcre2-8 -lintl -ltegra_hal -L/usr/local/lib/gstreamer-1.0 -lgstcoreelements -lgstapp -lgstplayback -lgstrawparse -lgsttcp -lgstvideoconvertscale -lgstvideotestsrc -lgstaudioparsers -lgstcodectimestamper  -lgstisomp4 -lgstopenh264 -lgstrtp -lgstrtpmanager -lgstrtsp -lgstudp -lgstvideoparsersbad -lmount -lblkid -leconf"

# actually do the build
FROM --platform=${TARGETPLATFORM} wasmvision-build-${TARGETARCH} AS wasmvision-build
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags "-linkmode 'external' -extldflags '-static'" -tags netgo,osusergo,customenv -o /build/wasmvision ./cmd/wasmvision


# final stage: create a minimal image with the wasmvision binary
FROM --platform=${TARGETPLATFORM} alpine:3.20 AS wasmvision-final

COPY --from=wasmvision-build /build/wasmvision /run/wasmvision

COPY ./processors/*.wasm /processors/

ENTRYPOINT ["/run/wasmvision"]
