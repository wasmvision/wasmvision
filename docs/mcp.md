# MCP Server

wasmVision includes experimental support for the [Model Context Protocol (MCP)](https://modelcontextprotocol.info/) by providing a [(MCP) Server](https://modelcontextprotocol.info/specification/draft/server/).

The MCP server implementation uses the "HTTP with SSE" transport to facilitate being used in various networking related scenarios.

## Using the MCP Server

To turn on the MCP server, use the `wasmvision run` command with the `--mcp-server=true` flag.

You can also set the port to be used for the MCP server with the `--mcp-port` flag. For example, `--mcp-port=http://192.168.1.13:1313`

## Resources

wasmVision provides [MCP resources](https://modelcontextprotocol.io/docs/concepts/resources) for image data.

### `images://input`

Reading from the `images://input` resource returns the current input image before the running processing pipeline. The returned image data is of MIME type `image/jpeg` and has been base64 encoded.

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "data": "base64-encoded-image-data-here"
    }
}
```

### `images://output`

Reading from the `images://output` resource returns the current output image after the running processors in the pipeline have processed it. The returned image data is of MIME type `image/jpeg` and has been base64 encoded.

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "data": "base64-encoded-image-data-here"
    }
}
```
