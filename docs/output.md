# wasmVision Output Destinations

There are several different ways you can output the results of the video that you process.

This can be either via streaming, or by saving to a file.

## MPJEG Stream

By default, wasmVision can stream the output in MJPEG format.

When using wasmVision, this is configured by using the `--output=mjpeg` flag.

To control the stream destination port, you can use the `--destination` (`-d` for short) flag like this:

```shell
wasmvision run -p hello.wasm -o mjpeg -d :8081
```

## Save to File

wasmVision can write data to video files in several different formats.

When using wasmVision, this is configured by using the `--output=file` flag.

To control the file type and destination, you can use the `--destination` (`-d` for short) flag like this:

```shell
wasmvision run -p hello.wasm -o file -d /path/to/video/filename.avi
```

## Other streaming formats

You can also output streams in a number of different formats using GStreamer. More info about this will go here soon.
