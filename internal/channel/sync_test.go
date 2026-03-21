package channel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverChannels(t *testing.T) {
	tmpDir := t.TempDir()
	projectsDir := filepath.Join(tmpDir, "projects")
	os.MkdirAll(projectsDir, 0755)
	os.MkdirAll(filepath.Join(projectsDir, "test-project-a"), 0755)
	os.MkdirAll(filepath.Join(projectsDir, "test-project-b"), 0755)
	os.WriteFile(filepath.Join(projectsDir, "somefile.txt"), []byte(""), 0644)
	discovered, err := DiscoverChannels(projectsDir)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(discovered) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(discovered))
	}
}

func TestDiffSync(t *testing.T) {
	existing := &Registry{
		Channels: []Channel{
			{ID: "a", Path: "/a", Name: "a", Status: StatusHealthy},
			{ID: "b", Path: "/b", Name: "b", Status: StatusHealthy},
			{ID: "c", Path: "/c", Name: "c", Pinned: true, Group: "work", Status: StatusHealthy},
		},
	}
	discovered := []Channel{
		{ID: "a", Path: "/a", Name: "a"},
		{ID: "c", Path: "/c", Name: "c"},
		{ID: "d", Path: "/d", Name: "d"},
	}
	result := DiffSync(existing, discovered)
	if len(result.Channels) != 3 {
		t.Fatalf("expected 3 channels, got %d", len(result.Channels))
	}
	ch := result.FindByID("c")
	if ch == nil {
		t.Fatal("expected channel c")
	}
	if !ch.Pinned {
		t.Error("expected c to be pinned (preserved)")
	}
	if ch.Group != "work" {
		t.Error("expected c group to be work (preserved)")
	}
	if result.FindByID("b") != nil {
		t.Error("expected channel b to be removed")
	}
	if result.FindByID("d") == nil {
		t.Error("expected channel d to be added")
	}
}

func TestDiscoverChannelsEmpty(t *testing.T) {
	channels, err := DiscoverChannels("/nonexistent/path")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(channels) != 0 {
		t.Errorf("expected 0 channels, got %d", len(channels))
	}
}
