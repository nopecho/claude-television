package claude

import (
	"path/filepath"
	"testing"
)

func TestExtractMCPServers(t *testing.T) {
	path := filepath.Join("testdata", "settings_with_mcp.json")
	settings, err := ParseSettings(path)
	if err != nil {
		t.Fatalf("parse settings: %v", err)
	}

	servers, err := ExtractMCPServers(settings, "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(servers))
	}

	var stdio, sse *MCPServer
	for i := range servers {
		if servers[i].Name == "my-server" {
			stdio = &servers[i]
		}
		if servers[i].Name == "remote-server" {
			sse = &servers[i]
		}
	}

	if stdio == nil || sse == nil {
		t.Fatal("expected both servers to be found")
	}
	if stdio.Type != "stdio" {
		t.Errorf("expected stdio type, got %s", stdio.Type)
	}
	if stdio.Command != "npx" {
		t.Errorf("expected npx command, got %s", stdio.Command)
	}
	if len(stdio.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(stdio.Args))
	}
	if stdio.Source != "global" {
		t.Errorf("expected global source, got %s", stdio.Source)
	}
	if sse.Type != "sse" {
		t.Errorf("expected sse type, got %s", sse.Type)
	}
	if sse.URL != "http://localhost:3000/sse" {
		t.Errorf("expected URL, got %s", sse.URL)
	}
}

func TestExtractMCPServersNil(t *testing.T) {
	servers, err := ExtractMCPServers(nil, "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(servers) != 0 {
		t.Errorf("expected 0 servers for nil settings, got %d", len(servers))
	}
}
