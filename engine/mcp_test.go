package engine

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/wasmvision/wasmvision/cv"
	"gocv.io/x/gocv"
)

func TestMCPServer(t *testing.T) {
	t.Run("new MCP server", func(t *testing.T) {
		port := ":8081"

		s := NewMCPServer(port)

		if s.Port != port {
			t.Errorf("unexpected port: %s", s.Port)
		}

		if s.outputFrames == nil {
			t.Errorf("unexpected nil frames")
		}
	})
}

func TestMCPServerStart(t *testing.T) {
	t.Run("start MCP server start", func(t *testing.T) {
		port := ":8081"

		s := NewMCPServer(port)

		s.Start()
		defer s.Close()
		img := gocv.IMRead("../images/wasmvision-logo.png", gocv.IMReadColor)
		frm := cv.NewFrame(img)
		if err := s.PublishOutput(frm); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestMCPServerEndpoint(t *testing.T) {
	t.Run("start MCP server start", func(t *testing.T) {
		port := "http://localhost:8081"

		s := NewMCPServer(port)
		s.Start()
		defer s.Close()

		sseResp, err := http.Get(fmt.Sprintf("%s/sse", port))
		if err != nil {
			t.Fatalf("Failed to connect to SSE endpoint: %v", err)
		}
		defer sseResp.Body.Close()

		buf := make([]byte, 1024)
		n, err := sseResp.Body.Read(buf)
		if err != nil {
			t.Fatalf("Failed to read SSE response body: %v", err)
		}

		if n == 0 {
			t.Fatalf("SSE response body is empty")
		}

		if !strings.Contains(string(buf[:n]), "event: endpoint") {
			t.Fatalf("SSE response body does not contain event: endpoint")
		}
	})
}
