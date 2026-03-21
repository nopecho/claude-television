package channel

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestCacheSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cache := NewCache(dir, 24*time.Hour)
	entry := &CacheEntry{
		ChannelID: "test-channel",
		Data: ChannelData{
			ClaudeMD: &claude.ClaudeMD{
				Path: "/test/CLAUDE.md", LineCount: 10,
				Sections: []string{"Build", "Test"},
			},
		},
		FileMtimes: map[string]time.Time{},
	}
	if err := cache.Save(entry); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := cache.Load("test-channel")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.ChannelID != "test-channel" {
		t.Errorf("unexpected ID: %s", loaded.ChannelID)
	}
	if loaded.Data.ClaudeMD == nil {
		t.Fatal("expected ClaudeMD in cached data")
	}
	if loaded.Data.ClaudeMD.LineCount != 10 {
		t.Errorf("expected 10 lines, got %d", loaded.Data.ClaudeMD.LineCount)
	}
}

func TestCacheIsValid(t *testing.T) {
	dir := t.TempDir()
	cache := NewCache(dir, 24*time.Hour)
	trackedFile := filepath.Join(dir, "tracked.json")
	os.WriteFile(trackedFile, []byte("{}"), 0644)
	info, _ := os.Stat(trackedFile)
	entry := &CacheEntry{
		ChannelID:  "test",
		Data:       ChannelData{},
		FileMtimes: map[string]time.Time{trackedFile: info.ModTime()},
	}
	cache.Save(entry)
	if !cache.IsValid("test", []string{trackedFile}) {
		t.Error("expected cache to be valid")
	}
	os.WriteFile(trackedFile, []byte(`{"changed": true}`), 0644)
	if cache.IsValid("test", []string{trackedFile}) {
		t.Error("expected cache to be invalid after file change")
	}
}

func TestCacheNewFileDetection(t *testing.T) {
	dir := t.TempDir()
	cache := NewCache(dir, 24*time.Hour)
	entry := &CacheEntry{
		ChannelID:  "test",
		Data:       ChannelData{},
		FileMtimes: map[string]time.Time{},
	}
	cache.Save(entry)
	newFile := filepath.Join(dir, "new.json")
	os.WriteFile(newFile, []byte("{}"), 0644)
	if cache.IsValid("test", []string{newFile}) {
		t.Error("expected cache to be invalid when new expected file appears")
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	dir := t.TempDir()
	cache := NewCache(dir, 1*time.Millisecond)
	entry := &CacheEntry{
		ChannelID:  "test",
		Data:       ChannelData{},
		FileMtimes: map[string]time.Time{},
	}
	cache.Save(entry)
	// Sleep to ensure TTL expires
	time.Sleep(5 * time.Millisecond)
	if cache.IsValid("test", nil) {
		t.Error("expected cache to be expired")
	}
}
