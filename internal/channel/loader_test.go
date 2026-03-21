package channel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadChannelData(t *testing.T) {
	projectDir := t.TempDir()
	claudeDir := filepath.Join(projectDir, ".claude")
	os.MkdirAll(claudeDir, 0755)

	settings := `{"model": "claude-sonnet-4-6", "permissions": {"allow": ["Read"]}}`
	os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(settings), 0644)

	claudeMD := "# Project\n\n## Build\n\ngo build\n\n## Test\n\ngo test\n"
	os.WriteFile(filepath.Join(projectDir, "CLAUDE.md"), []byte(claudeMD), 0644)

	ch := &Channel{ID: "test", Path: projectDir, Name: "test"}

	data, mtimes, err := LoadChannelData(ch, "")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if data.Settings == nil {
		t.Fatal("expected settings")
	}
	if data.Settings.Model != "claude-sonnet-4-6" {
		t.Errorf("expected claude-sonnet-4-6, got %s", data.Settings.Model)
	}
	if data.ClaudeMD == nil {
		t.Fatal("expected CLAUDE.md")
	}
	if len(data.ClaudeMD.Sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(data.ClaudeMD.Sections))
	}
	if len(mtimes) == 0 {
		t.Error("expected file mtimes to be collected")
	}
}

func TestLoadChannelDataEmpty(t *testing.T) {
	projectDir := t.TempDir()
	ch := &Channel{ID: "empty", Path: projectDir, Name: "empty"}

	data, _, err := LoadChannelData(ch, "")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if data.Settings != nil {
		t.Error("expected nil settings for empty project")
	}
	if data.ClaudeMD != nil {
		t.Error("expected nil ClaudeMD for empty project")
	}
}

func TestExpectedFiles(t *testing.T) {
	ch := &Channel{ID: "test", Path: "/projects/test"}
	files := ExpectedFiles(ch)
	if len(files) != 3 {
		t.Fatalf("expected 3 expected files, got %d", len(files))
	}
}
