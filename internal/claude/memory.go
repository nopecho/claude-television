package claude

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MemoryFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func ScanMemoryFiles(dir string) ([]MemoryFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan memory: %w", err)
	}
	var result []MemoryFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if strings.ToUpper(e.Name()) == "MEMORY.MD" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		mf := parseMemoryFrontmatter(path)
		mf.Path = path
		if mf.Name == "" {
			mf.Name = strings.TrimSuffix(e.Name(), ".md")
		}
		result = append(result, mf)
	}
	return result, nil
}

func parseMemoryFrontmatter(path string) MemoryFile {
	f, err := os.Open(path)
	if err != nil {
		return MemoryFile{}
	}
	defer f.Close()
	var mf MemoryFile
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter {
			if k, v, ok := parseYAMLLine(trimmed); ok {
				switch k {
				case "name":
					mf.Name = v
				case "description":
					mf.Description = v
				case "type":
					mf.Type = v
				}
			}
		}
	}
	return mf
}

func parseYAMLLine(line string) (key, value string, ok bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	return key, value, true
}
