package channel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	ChannelID  string               `json:"channel_id"`
	Data       ChannelData          `json:"data"`
	FileMtimes map[string]time.Time `json:"file_mtimes"`
	CachedAt   time.Time            `json:"cached_at"`
}

type Cache struct {
	dir string
	ttl time.Duration
}

func NewCache(dir string, ttl time.Duration) *Cache {
	return &Cache{dir: dir, ttl: ttl}
}

func (c *Cache) Save(entry *CacheEntry) error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	entry.CachedAt = time.Now()
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}
	path := filepath.Join(c.dir, entry.ChannelID+".json")
	return os.WriteFile(path, data, 0644)
}

func (c *Cache) Load(channelID string) (*CacheEntry, error) {
	path := filepath.Join(c.dir, channelID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read cache: %w", err)
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("parse cache: %w", err)
	}
	return &entry, nil
}

func (c *Cache) LoadIfValid(channelID string, expectedFiles []string) (*CacheEntry, bool) {
	entry, err := c.Load(channelID)
	if err != nil || entry == nil {
		return nil, false
	}
	if time.Since(entry.CachedAt) > c.ttl {
		return nil, false
	}
	for _, path := range expectedFiles {
		if _, tracked := entry.FileMtimes[path]; !tracked {
			if _, err := os.Stat(path); err == nil {
				return nil, false
			}
		}
	}
	for path, cachedMtime := range entry.FileMtimes {
		info, err := os.Stat(path)
		if err != nil || info.ModTime().After(cachedMtime) {
			return nil, false
		}
	}
	return entry, true
}

func (c *Cache) IsValid(channelID string, expectedFiles []string) bool {
	_, valid := c.LoadIfValid(channelID, expectedFiles)
	return valid
}

func (c *Cache) Delete(channelID string) error {
	path := filepath.Join(c.dir, channelID+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
