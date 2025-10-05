package engine

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/wasmvision/wasmvision"
	"github.com/wasmvision/wasmvision/cv"
	"gocv.io/x/gocv"
)

func TestMCPServer(t *testing.T) {
	t.Run("new MCP server", func(t *testing.T) {
		port := ":8081"

		s := NewMCPServer(port)
		defer func() {
			s.Close()
			time.Sleep(500 * time.Millisecond)
		}()

		if s.Port != port {
			t.Errorf("unexpected port: %s", s.Port)
		}

		if s.outputFrames == nil {
			t.Errorf("unexpected nil frames")
		}
	})
}

func TestMCPServerPublishFrame(t *testing.T) {
	t.Run("start MCP server start", func(t *testing.T) {
		port := ":8081"

		s := NewMCPServer(port)

		s.mcpServer = server.NewMCPServer("wasmvision-test", wasmvision.Version())
		s.AddImageInputResource()
		s.AddImageOutputResource()
		s.AddProcessorDatastoreResource()

		s.httpServer, _ = NewTestStreamableHTTPServer(s.mcpServer, server.WithEndpointPath("/mcp"))

		s.StartPublishing()

		defer func() {
			s.Close()
			time.Sleep(500 * time.Millisecond)
		}()

		img := gocv.IMRead("../images/wasmvision-logo.png", gocv.IMReadColor)
		frm := cv.NewFrame(img)
		if err := s.PublishOutput(frm); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestMCPServerEndpoint(t *testing.T) {
	t.Run("start MCP server start", func(t *testing.T) {
		port := ":8081"

		s := NewMCPServer(port)

		s.mcpServer = server.NewMCPServer("wasmvision-test", wasmvision.Version())
		s.AddImageInputResource()
		s.AddImageOutputResource()
		s.AddProcessorDatastoreResource()

		httpSrv, srv := NewTestStreamableHTTPServer(s.mcpServer, server.WithEndpointPath("/mcp"))
		s.httpServer = httpSrv
		s.Port = srv.URL

		defer func() {
			s.Close()
			time.Sleep(500 * time.Millisecond)
		}()

		httpResp, err := http.Get(srv.URL + "/mcp")
		if err != nil {
			t.Fatalf("Failed to connect to MCP endpoint: %v", err)
		}
		defer httpResp.Body.Close()
		if httpResp.StatusCode != http.StatusOK {
			t.Fatalf("MCP response body is empty")
		}
	})
}

// NewTestStreamableHTTPServer creates a test server for testing purposes
func NewTestStreamableHTTPServer(s *server.MCPServer, opts ...server.StreamableHTTPOption) (*server.StreamableHTTPServer, *httptest.Server) {
	httpServer := server.NewStreamableHTTPServer(s, opts...)
	testServer := httptest.NewServer(httpServer)
	return httpServer, testServer
}
