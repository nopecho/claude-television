package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestScanProjectsMeta(t *testing.T) {
	dir := t.TempDir()

	// Project with memory and sessions
	proj1 := filepath.Join(dir, "-Users-test-myproject")
	os.MkdirAll(filepath.Join(proj1, "memory"), 0o755)
	os.MkdirAll(filepath.Join(proj1, "session-abc"), 0o755)

	// Project with only memory
	proj2 := filepath.Join(dir, "-Users-test-other")
	os.MkdirAll(filepath.Join(proj2, "memory"), 0o755)

	// A file (should be ignored)
	os.WriteFile(filepath.Join(dir, "some-file.txt"), []byte("hi"), 0o644)

	projects, err := claude.ScanProjectsMeta(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("projects len = %d, want 2", len(projects))
	}

	byName := make(map[string]claude.ProjectMeta)
	for _, p := range projects {
		byName[p.EncodedName] = p
	}

	p1 := byName["-Users-test-myproject"]
	if !p1.HasMemory {
		t.Error("expected HasMemory = true")
	}
	if !p1.HasSessions {
		t.Error("expected HasSessions = true")
	}

	p2 := byName["-Users-test-other"]
	if !p2.HasMemory {
		t.Error("expected HasMemory = true")
	}
	if p2.HasSessions {
		t.Error("expected HasSessions = false")
	}
}

func TestScanProjectsMeta_NotExist(t *testing.T) {
	projects, err := claude.ScanProjectsMeta("/nonexistent/path/projects")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if projects != nil {
		t.Errorf("expected nil, got %v", projects)
	}
}
