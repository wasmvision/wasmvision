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
	mcpServer          *server.MCPServer
	sseServer          *server.SSEServer
	Port               string
	inputFrames        chan *cv.Frame
	currentInputFrame  gocv.Mat
	inputFrameMut      sync.Mutex
	outputFrames       chan *cv.Frame
	currentOutputFrame gocv.Mat
	outputFrameMut     sync.Mutex
}

// NewMCPServer creates a new MCPServer instance with the given port.
func NewMCPServer(port string) *MCPServer {
	return &MCPServer{
		Port:               port,
		inputFrames:        make(chan *cv.Frame, framebufferSize),
		currentInputFrame:  gocv.NewMat(),
		outputFrames:       make(chan *cv.Frame, framebufferSize),
		currentOutputFrame: gocv.NewMat(),
	}
}

var (
	imagesInputResource = mcp.NewResource(
		"images://input",
		"input",
		mcp.WithResourceDescription("Current input image frame"),
		mcp.WithMIMEType("image/jpeg"),
	)

	imagesOutputResource = mcp.NewResource(
		"images://output",
		"output",
		mcp.WithResourceDescription("Current output image frame"),
		mcp.WithMIMEType("image/jpeg"),
	)
)

// Start starts the NewMCPServer server.
func (s *MCPServer) Start() error {
	s.mcpServer = server.NewMCPServer("wasmvision", wasmvision.Version(),
		server.WithResourceCapabilities(false, false))
	s.sseServer = server.NewSSEServer(s.mcpServer,
		server.WithBaseURL(getURL(s.Port)),
	)

	s.mcpServer.AddResource(imagesInputResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		resource, err := handleImageResource(&s.inputFrameMut, &s.currentInputFrame, imagesInputResource.URI)
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			resource,
		}, nil
	})

	s.mcpServer.AddResource(imagesOutputResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		resource, err := handleImageResource(&s.outputFrameMut, &s.currentOutputFrame, imagesOutputResource.URI)
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			resource,
		}, nil
	})

	go s.sseServer.Start(getPort(s.Port))
	go s.publishInputFrames()
	go s.publishOutputFrames()

	return nil
}

// Close closes the MCPServer server.
func (s *MCPServer) Close() {
	close(s.outputFrames)
	if s.sseServer != nil {
		s.sseServer.Shutdown(context.Background())
	}
}

// Publish publishes an input frame to the MCP server.
func (s *MCPServer) PublishInput(frm *cv.Frame) error {
	s.inputFrames <- frm
	return nil
}

func (s *MCPServer) publishInputFrames() {
	for frame := range s.inputFrames {
		// if there is outstanding request for frame resource, update it
		s.inputFrameMut.Lock()
		frame.Image.CopyTo(&s.currentInputFrame)
		s.inputFrameMut.Unlock()
	}
}

// Publish publishes an output frame to the MCP server.
func (s *MCPServer) PublishOutput(frm *cv.Frame) error {
	s.outputFrames <- frm
	return nil
}

func (s *MCPServer) publishOutputFrames() {
	for frame := range s.outputFrames {
		// if there is outstanding request for frame resource, update it
		s.outputFrameMut.Lock()
		frame.Image.CopyTo(&s.currentOutputFrame)
		s.outputFrameMut.Unlock()
	}
}

func handleImageResource(mut *sync.Mutex, frame *gocv.Mat, uri string) (mcp.BlobResourceContents, error) {
	mut.Lock()
	defer mut.Unlock()

	// Create a copy of the current frame for encoding
	frameCopy := gocv.NewMat()
	defer frameCopy.Close()
	frame.CopyTo(&frameCopy)

	buf, err := gocv.IMEncode(".jpg", frameCopy)
	if err != nil {
		slog.Error(fmt.Sprintf("error encoding frame: %v", err))
		return mcp.BlobResourceContents{}, err
	}
	defer buf.Close()

	encodedImage := base64.StdEncoding.EncodeToString(buf.GetBytes())

	return mcp.BlobResourceContents{
		URI:      uri,
		MIMEType: "image/jpeg",
		Blob:     encodedImage,
	}, nil
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
