package channel

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRegistrySaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	reg := &Registry{
		Channels: []Channel{
			{
				ID:         "-Users-test-project",
				Path:       "/Users/test/project",
				Name:       "project",
				Status:     StatusHealthy,
				LastSynced: time.Now().Truncate(time.Second),
			},
		},
		UpdatedAt: time.Now().Truncate(time.Second),
	}
	path := filepath.Join(dir, "channels.json")
	if err := SaveRegistry(reg, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRegistry(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Channels) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(loaded.Channels))
	}
	if loaded.Channels[0].ID != "-Users-test-project" {
		t.Errorf("unexpected ID: %s", loaded.Channels[0].ID)
	}
	if loaded.Channels[0].Status != StatusHealthy {
		t.Errorf("unexpected status: %s", loaded.Channels[0].Status)
	}
}

func TestLoadRegistryNotFound(t *testing.T) {
	reg, err := LoadRegistry("/nonexistent/channels.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(reg.Channels) != 0 {
		t.Errorf("expected empty registry, got %d channels", len(reg.Channels))
	}
}
