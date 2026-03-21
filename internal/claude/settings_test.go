package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestParseSettings(t *testing.T) {
	path := filepath.Join("testdata", "settings.json")
	s, err := claude.ParseSettings(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Model != "opus" {
		t.Errorf("model = %q, want %q", s.Model, "opus")
	}
	if s.Language != "korean" {
		t.Errorf("language = %q, want %q", s.Language, "korean")
	}
	if got := s.Env["CLAUDE_CODE_SHELL"]; got != "zsh" {
		t.Errorf("env CLAUDE_CODE_SHELL = %q, want %q", got, "zsh")
	}
	if len(s.Permissions.Allow) != 2 {
		t.Errorf("permissions.allow len = %d, want 2", len(s.Permissions.Allow))
	}
	if s.Permissions.Allow[0] != "Bash(go:*)" {
		t.Errorf("permissions.allow[0] = %q, want %q", s.Permissions.Allow[0], "Bash(go:*)")
	}
	if len(s.Permissions.Deny) != 1 {
		t.Errorf("permissions.deny len = %d, want 1", len(s.Permissions.Deny))
	}
	if !s.EnabledPlugins["superpowers@claude-plugins-official"] {
		t.Error("expected superpowers plugin to be enabled")
	}
	if s.EnabledPlugins["obsidian@obsidian-skills"] {
		t.Error("expected obsidian plugin to be disabled")
	}
	if s.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestParseSettings_NotFound(t *testing.T) {
	_, err := claude.ParseSettings("testdata/nonexistent.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
