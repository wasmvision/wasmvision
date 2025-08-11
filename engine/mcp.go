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
	"github.com/wasmvision/wasmvision/datastore"
	"gocv.io/x/gocv"
)

// MCPServer represents a MCP server currently providing a resource of the frames
// being processed by wasmVision.
type MCPServer struct {
	mcpServer          *server.MCPServer
	httpServer         *server.StreamableHTTPServer
	Port               string
	inputFrames        chan *cv.Frame
	currentInputFrame  gocv.Mat
	inputFrameMut      sync.Mutex
	outputFrames       chan *cv.Frame
	currentOutputFrame gocv.Mat
	outputFrameMut     sync.Mutex
	ProcessorStore     *datastore.Processors
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

// Init inits the NewMCPServer server.
func (s *MCPServer) Init() error {
	s.mcpServer = server.NewMCPServer("wasmvision", wasmvision.Version(), server.WithResourceCapabilities(false, false), server.WithLogging())
	s.httpServer = server.NewStreamableHTTPServer(s.mcpServer)

	s.AddImageInputResource()
	s.AddImageOutputResource()
	s.AddProcessorDatastoreResource()

	return nil
}

// Start starts the NewMCPServer server.
func (s *MCPServer) Start() error {
	slog.Warn("MCP server starting", "port", s.Port)

	go s.httpServer.Start(getPort(s.Port))
	s.StartPublishing()

	return nil
}

// AddImageInputResource adds the image input resource.
func (s *MCPServer) AddImageInputResource() error {
	imagesInputResource := mcp.NewResource(
		"images://input",
		"input",
		mcp.WithResourceDescription("Current input image frame"),
		mcp.WithMIMEType("image/jpeg"),
	)

	s.mcpServer.AddResource(imagesInputResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		resource, err := handleImageResource(&s.inputFrameMut, &s.currentInputFrame, imagesInputResource.URI)
		if err != nil {
			slog.Error("unable to handle MCP input resource", "error", err.Error())
			return nil, err
		}

		return []mcp.ResourceContents{
			resource,
		}, nil
	})

	return nil
}

// AddImageOutputResource adds the image output resource.
func (s *MCPServer) AddImageOutputResource() error {
	imagesOutputResource := mcp.NewResource(
		"images://output",
		"output",
		mcp.WithResourceDescription("Current output image frame"),
		mcp.WithMIMEType("image/jpeg"),
	)

	s.mcpServer.AddResource(imagesOutputResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		resource, err := handleImageResource(&s.outputFrameMut, &s.currentOutputFrame, imagesOutputResource.URI)
		if err != nil {
			slog.Error("unable to handle MCP output resource", "error", err.Error())
			return nil, err
		}

		return []mcp.ResourceContents{
			resource,
		}, nil
	})

	return nil
}

// AddProcessorDatastoreResource adds the Processor datastore resource.
func (s *MCPServer) AddProcessorDatastoreResource() error {
	datastoreResource := mcp.NewResource(
		"data://processor/{processor}/{key}",
		"datastore data for processor",
		mcp.WithResourceDescription("processor data from datastore for a specific processor and key"),
		mcp.WithMIMEType("application/json"),
	)

	s.mcpServer.AddResource(datastoreResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		var processor, key, content string
		if _, err := fmt.Sscanf(datastoreResource.URI, "data://processor/%s/%s", &processor, &key); err != nil {
			return nil, fmt.Errorf("invalid resource URI format: %w", err)
		}

		var ok bool
		content, ok = s.ProcessorStore.Get(processor, key)
		if !ok {
			content = "{\"error\": \"could not find processor datastore information\"}"
		}

		resource := mcp.TextResourceContents{
			URI:      datastoreResource.URI,
			MIMEType: datastoreResource.MIMEType,
			Text:     content,
		}

		return []mcp.ResourceContents{
			resource,
		}, nil
	})

	return nil
}

// StartPublishing starts publishing frames.
func (s *MCPServer) StartPublishing() error {
	go s.publishInputFrames()
	go s.publishOutputFrames()

	return nil
}

// Close closes the MCPServer server.
func (s *MCPServer) Close() {
	close(s.outputFrames)
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			slog.Error("problem shutting down MCP server", "error", err)
		}
	}
}

// Publish publishes an input frame to the MCP server.
func (s *MCPServer) PublishInput(frm *cv.Frame) error {
	s.inputFrames <- frm
	return nil
}

func (s *MCPServer) publishInputFrames() {
	for frame := range s.inputFrames {
		if frame == nil || frame.Empty() || frame.Image.Ptr() == nil || frame.Image.Empty() {
			slog.Error("empty frame")
			continue
		}

		s.inputFrameMut.Lock()
		frame.Image.CopyTo(&s.currentInputFrame)
		s.inputFrameMut.Unlock()
		frame.Close()
	}
}

// Publish publishes an output frame to the MCP server.
func (s *MCPServer) PublishOutput(frm *cv.Frame) error {
	s.outputFrames <- frm
	return nil
}

func (s *MCPServer) publishOutputFrames() {
	for frame := range s.outputFrames {
		if frame == nil || frame.Empty() || frame.Image.Ptr() == nil || frame.Image.Empty() {
			slog.Error("empty frame")
			continue
		}

		s.outputFrameMut.Lock()
		frame.Image.CopyTo(&s.currentOutputFrame)
		s.outputFrameMut.Unlock()
		frame.Close()
	}
}

func handleImageResource(mut *sync.Mutex, frame *gocv.Mat, uri string) (mcp.BlobResourceContents, error) {
	mut.Lock()
	defer mut.Unlock()

	buf, err := gocv.IMEncode(".jpg", *frame)
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

func getPort(port string) string {
	if port == "" {
		return ":9090"
	}
	if port[0] == ':' {
		return port
	}

	u, err := url.Parse(port)
	if err != nil {
		return ":9090"
	}

	_, p, _ := net.SplitHostPort(u.Host)
	if p == "" {
		return ":9090"
	}
	return ":" + p
}
