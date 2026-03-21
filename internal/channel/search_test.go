package channel

import (
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestFuzzySearch(t *testing.T) {
	channels := []Channel{
		{ID: "1", Name: "my-api-server", Path: "/projects/my-api-server", Group: "work"},
		{ID: "2", Name: "web-frontend", Path: "/projects/web-frontend", Group: "work"},
		{ID: "3", Name: "personal-blog", Path: "/home/blog", Group: "personal"},
		{ID: "4", Name: "api-gateway", Path: "/projects/api-gateway"},
	}

	results := FuzzySearch(channels, "api")
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'api', got %d", len(results))
	}

	results = FuzzySearch(channels, "web")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'web', got %d", len(results))
	}

	results = FuzzySearch(channels, "")
	if len(results) != 4 {
		t.Errorf("expected all 4 for empty query, got %d", len(results))
	}
}

func TestContentSearch(t *testing.T) {
	channels := []Channel{
		{
			ID: "1", Name: "project-a",
			Data: &ChannelData{
				MCPServers: []claude.MCPServer{{Name: "context7"}},
				Settings:   &claude.Settings{Model: "opus"},
			},
		},
		{
			ID: "2", Name: "project-b",
			Data: &ChannelData{
				Hooks: []claude.HookDetail{{Event: "pre-commit", Command: "lint"}},
			},
		},
		{
			ID: "3", Name: "project-c",
			Data: &ChannelData{},
		},
	}

	tests := []struct {
		query    string
		wantIDs  []string
	}{
		{"context7", []string{"1"}},
		{"opus", []string{"1"}},
		{"pre-commit", []string{"2"}},
		{"nonexistent", nil},
		{"", []string{"1", "2", "3"}},
	}

	for _, tt := range tests {
		results := ContentSearch(channels, tt.query)
		if tt.wantIDs == nil {
			if len(results) != 0 {
				t.Errorf("query=%q: expected 0 results, got %d", tt.query, len(results))
			}
			continue
		}
		if len(results) != len(tt.wantIDs) {
			t.Errorf("query=%q: expected %d results, got %d", tt.query, len(tt.wantIDs), len(results))
			continue
		}
		for i, r := range results {
			if r.ID != tt.wantIDs[i] {
				t.Errorf("query=%q: result[%d].ID = %s, want %s", tt.query, i, r.ID, tt.wantIDs[i])
			}
		}
	}
}
