package channel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nopecho/claude-television/internal/util"
)

func DiscoverChannels(projectsDir string) ([]Channel, error) {
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read projects dir: %w", err)
	}
	var channels []Channel
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id := e.Name()
		decoded := util.DecodeProjectPath(id)
		name := filepath.Base(decoded)
		status := StatusHealthy
		if _, err := os.Stat(decoded); err != nil {
			status = StatusError
		}
		channels = append(channels, Channel{
			ID: id, Path: decoded, Name: name,
			Status: status, LastSynced: time.Now(),
		})
	}
	return channels, nil
}

func DiffSync(existing *Registry, discovered []Channel) *Registry {
	result := &Registry{UpdatedAt: time.Now()}
	existingMap := map[string]*Channel{}
	for i := range existing.Channels {
		existingMap[existing.Channels[i].ID] = &existing.Channels[i]
	}
	for _, d := range discovered {
		if prev, ok := existingMap[d.ID]; ok {
			d.Pinned = prev.Pinned
			d.Group = prev.Group
		}
		result.Channels = append(result.Channels, d)
	}
	return result
}

