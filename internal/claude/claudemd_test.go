package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestParseClaudeMD(t *testing.T) {
	path := filepath.Join("testdata", "CLAUDE.md")
	md, err := claude.ParseClaudeMD(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if md.Path != path {
		t.Errorf("path = %q, want %q", md.Path, path)
	}
	if md.LineCount != 10 {
		t.Errorf("lineCount = %d, want 10", md.LineCount)
	}
	expectedSections := []string{"Build", "Testing", "Conventions"}
	if len(md.Sections) != len(expectedSections) {
		t.Fatalf("sections len = %d, want %d", len(md.Sections), len(expectedSections))
	}
	for i, s := range expectedSections {
		if md.Sections[i] != s {
			t.Errorf("sections[%d] = %q, want %q", i, md.Sections[i], s)
		}
	}
	if md.Content == "" {
		t.Error("expected content to be non-empty")
	}
}

func TestParseClaudeMD_NotFound(t *testing.T) {
	_, err := claude.ParseClaudeMD("testdata/nonexistent.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
