package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestScanLocalSkills(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "my-skill"), 0o755)
	os.MkdirAll(filepath.Join(dir, "another-skill"), 0o755)
	// Create a file (should be ignored)
	os.WriteFile(filepath.Join(dir, "not-a-skill.txt"), []byte("hi"), 0o644)

	skills, err := claude.ScanLocalSkills(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("skills len = %d, want 2", len(skills))
	}

	names := make(map[string]bool)
	for _, s := range skills {
		names[s.Name] = true
	}
	if !names["my-skill"] {
		t.Error("expected my-skill")
	}
	if !names["another-skill"] {
		t.Error("expected another-skill")
	}
}

func TestScanLocalSkills_NotExist(t *testing.T) {
	skills, err := claude.ScanLocalSkills("/nonexistent/path/skills")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skills != nil {
		t.Errorf("expected nil, got %v", skills)
	}
}
