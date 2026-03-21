package claude

import (
	"path/filepath"
	"testing"
)

func TestScanMemoryFiles(t *testing.T) {
	dir := filepath.Join("testdata", "memory")
	files, err := ScanMemoryFiles(dir)
	if err != nil {
		t.Fatalf("scan memory: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 memory files, got %d", len(files))
	}
	found := map[string]MemoryFile{}
	for _, f := range files {
		found[f.Name] = f
	}
	user, ok := found["user_role"]
	if !ok {
		t.Fatal("expected user_role memory file")
	}
	if user.Type != "user" {
		t.Errorf("expected type user, got %s", user.Type)
	}
	if user.Description != "User is a backend engineer" {
		t.Errorf("unexpected description: %s", user.Description)
	}
}

func TestScanMemoryFilesEmpty(t *testing.T) {
	files, err := ScanMemoryFiles("/nonexistent/path")
	if err != nil {
		t.Fatalf("expected nil error for nonexistent, got %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}
