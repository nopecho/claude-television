package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanProjects(t *testing.T) {
	root := t.TempDir()

	// project-a: has .claude/ dir with settings.json
	projA := filepath.Join(root, "project-a")
	os.MkdirAll(filepath.Join(projA, ".claude"), 0o755)
	os.WriteFile(filepath.Join(projA, ".claude", "settings.json"), []byte("{}"), 0o644)

	// project-b: has CLAUDE.md only
	projB := filepath.Join(root, "project-b")
	os.MkdirAll(projB, 0o755)
	os.WriteFile(filepath.Join(projB, "CLAUDE.md"), []byte("# hi"), 0o644)

	// project-c: nothing
	os.MkdirAll(filepath.Join(root, "project-c"), 0o755)

	// node_modules: should be ignored
	os.MkdirAll(filepath.Join(root, "node_modules"), 0o755)

	projects, err := ScanProjects([]string{root}, []string{"node_modules"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 3 {
		t.Fatalf("expected 3 projects, got %d", len(projects))
	}

	byName := make(map[string]Project)
	for _, p := range projects {
		byName[p.Name] = p
	}

	a := byName["project-a"]
	if !a.HasClaudeDir || a.HasClaudeMD || !a.HasSettings {
		t.Errorf("project-a flags wrong: %+v", a)
	}

	b := byName["project-b"]
	if b.HasClaudeDir || !b.HasClaudeMD || b.HasSettings {
		t.Errorf("project-b flags wrong: %+v", b)
	}

	c := byName["project-c"]
	if c.HasClaudeDir || c.HasClaudeMD || c.HasSettings {
		t.Errorf("project-c flags wrong: %+v", c)
	}
}

func TestScanProjects_NotExist(t *testing.T) {
	projects, err := ScanProjects([]string{"/tmp/nonexistent-ctv-path"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected 0 projects, got %d", len(projects))
	}
}

func TestScanProjects_IgnoresHidden(t *testing.T) {
	root := t.TempDir()

	os.MkdirAll(filepath.Join(root, ".hidden-project"), 0o755)
	os.MkdirAll(filepath.Join(root, "visible-project"), 0o755)

	projects, err := ScanProjects([]string{root}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "visible-project" {
		t.Errorf("expected visible-project, got %s", projects[0].Name)
	}
}
