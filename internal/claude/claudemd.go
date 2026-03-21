package claude

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ClaudeMD struct {
	Path      string   `json:"path"`
	LineCount int      `json:"line_count"`
	Sections  []string `json:"sections"`
	Content   string   `json:"content"`
}

func ParseClaudeMD(path string) (*ClaudeMD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open CLAUDE.md: %w", err)
	}
	defer f.Close()
	var (
		lines    int
		sections []string
		content  strings.Builder
	)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		content.WriteString(line)
		content.WriteString("\n")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			sections = append(sections, strings.TrimPrefix(trimmed, "## "))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan CLAUDE.md: %w", err)
	}
	return &ClaudeMD{Path: path, LineCount: lines, Sections: sections, Content: content.String()}, nil
}
