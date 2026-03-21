package channel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		decoded := decodeProjectPath(id)
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

func decodeProjectPath(encoded string) string {
	if !strings.HasPrefix(encoded, "-") {
		return encoded
	}
	parts := strings.Split(strings.TrimPrefix(encoded, "-"), "-")
	return bestEffortDecode(parts)
}

func bestEffortDecode(parts []string) string {
	if len(parts) == 0 {
		return "/"
	}
	current := "/"
	i := 0
	for i < len(parts) {
		found := false
		for j := len(parts); j > i; j-- {
			candidate := current + strings.Join(parts[i:j], "-")
			if _, err := os.Stat(candidate); err == nil {
				current = candidate + "/"
				i = j
				found = true
				break
			}
		}
		if !found {
			current += parts[i] + "/"
			i++
		}
	}
	return strings.TrimSuffix(current, "/")
}
