package claude

import (
	"fmt"
	"os"
	"path/filepath"
)

type Skill struct {
	Name string
	Path string
}

func ScanLocalSkills(skillsDir string) ([]Skill, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan skills: %w", err)
	}
	var skills []Skill
	for _, e := range entries {
		if e.IsDir() {
			skills = append(skills, Skill{Name: e.Name(), Path: filepath.Join(skillsDir, e.Name())})
		}
	}
	return skills, nil
}
