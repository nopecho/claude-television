package claude

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectMeta struct {
	EncodedName string
	DecodedPath string
	HasMemory   bool
	HasSessions bool
}

func ScanProjectsMeta(projectsDir string) ([]ProjectMeta, error) {
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan projects: %w", err)
	}
	var result []ProjectMeta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		dirPath := filepath.Join(projectsDir, name)
		meta := ProjectMeta{EncodedName: name, DecodedPath: decodeProjectPath(name)}
		if _, err := os.Stat(filepath.Join(dirPath, "memory")); err == nil {
			meta.HasMemory = true
		}
		subEntries, _ := os.ReadDir(dirPath)
		for _, se := range subEntries {
			if se.IsDir() && se.Name() != "memory" {
				meta.HasSessions = true
				break
			}
		}
		result = append(result, meta)
	}
	return result, nil
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
			if pathExists(candidate) {
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
