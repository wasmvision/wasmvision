package engine

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wasmvision/wasmvision"
	"github.com/wasmvision/wasmvision/cv"
	"gocv.io/x/gocv"
)

// MCPServer represents a MCP server currently providing a resource of the frames
// being processed by wasmVision.
type MCPServer struct {
	mcpServer    *server.MCPServer
	sseServer    *server.SSEServer
	Port         string
	frames       chan *cv.Frame
	currentFrame gocv.Mat
	frameMut     sync.Mutex
}

// NewMCPServer creates a new MCPServer instance with the given port.
func NewMCPServer(port string) *MCPServer {
	return &MCPServer{
		Port:         port,
		frames:       make(chan *cv.Frame, framebufferSize),
		currentFrame: gocv.NewMat(),
	}
}

// Start starts the NewMCPServer server.
func (s *MCPServer) Start() error {
	s.mcpServer = server.NewMCPServer("wasmvision", wasmvision.Version(),
		server.WithResourceCapabilities(false, false))
	s.sseServer = server.NewSSEServer(s.mcpServer,
		server.WithBaseURL(getURL(s.Port)),
	)

	resource := mcp.NewResource(
		"images://output",
		"output",
		mcp.WithResourceDescription("Current output image frame"),
		mcp.WithMIMEType("image/jpeg"),
	)

	// Add resource with its handler
	s.mcpServer.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		s.frameMut.Lock()
		defer s.frameMut.Unlock()

		// Create a copy of the current frame for encoding
		frameCopy := gocv.NewMat()
		defer frameCopy.Close()
		s.currentFrame.CopyTo(&frameCopy)

		buf, err := gocv.IMEncode(".jpg", frameCopy)
		if err != nil {
			slog.Error(fmt.Sprintf("error encoding frame: %v", err))
			return nil, err
		}
		defer buf.Close()

		encodedImage := base64.StdEncoding.EncodeToString(buf.GetBytes())

		return []mcp.ResourceContents{
			mcp.BlobResourceContents{
				URI:      "images://output",
				MIMEType: "image/jpeg",
				Blob:     encodedImage,
			},
		}, nil
	})

	go s.sseServer.Start(getPort(s.Port))
	go s.publishFrames()

	return nil
}

// Close closes the MCPServer server.
func (s *MCPServer) Close() {
	close(s.frames)
	if s.sseServer != nil {
		s.sseServer.Shutdown(context.Background())
	}
}

// Publish publishes a frame to the MJPEG stream.
func (s *MCPServer) Publish(frm *cv.Frame) error {
	s.frames <- frm
	return nil
}

func (s *MCPServer) publishFrames() {
	for frame := range s.frames {
		// if there is outstanding request for frame resource, update it
		s.frameMut.Lock()
		frame.Image.CopyTo(&s.currentFrame)
		s.frameMut.Unlock()
	}
}

func getURL(port string) string {
	if port == "" {
		return "http://localhost:5001"
	}
	if port[0] == ':' {
		return fmt.Sprintf("http://localhost%s", port)
	}
	return port
}

func getPort(port string) string {
	if port == "" {
		return ":5001"
	}
	if port[0] == ':' {
		return port
	}

	u, err := url.Parse(port)
	if err != nil {
		return ":5001"
	}

	_, p, _ := net.SplitHostPort(u.Host)
	if p == "" {
		return ":5001"
	}
	return ":" + p
}
