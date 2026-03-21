package claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nopecho/claude-television/internal/util"
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
		meta := ProjectMeta{EncodedName: name, DecodedPath: util.DecodeProjectPath(name)}
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
