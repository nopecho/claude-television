package claude_test

import (
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestExtractHooks(t *testing.T) {
	settings := &claude.Settings{
		Hooks: map[string][]claude.HookRule{
			"PreToolUse": {
				{
					Matcher: "Bash",
					Hooks: []claude.HookAction{
						{Type: "command", Command: "echo pre-hook", Async: false, Timeout: 5},
					},
				},
			},
			"PostToolUse": {
				{
					Matcher: "*",
					Hooks: []claude.HookAction{
						{Type: "command", Command: "echo post-hook", Async: true, Timeout: 10},
					},
				},
			},
		},
	}

	hooks, err := claude.ExtractHooks(settings, "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hooks) != 2 {
		t.Fatalf("hooks len = %d, want 2", len(hooks))
	}

	byEvent := make(map[string]claude.HookDetail)
	for _, h := range hooks {
		byEvent[h.Event] = h
	}

	pre := byEvent["PreToolUse"]
	if pre.Matcher != "Bash" {
		t.Errorf("matcher = %q, want %q", pre.Matcher, "Bash")
	}
	if pre.Command != "echo pre-hook" {
		t.Errorf("command = %q, want %q", pre.Command, "echo pre-hook")
	}
	if pre.Source != "global" {
		t.Errorf("source = %q, want %q", pre.Source, "global")
	}
	if pre.Async {
		t.Error("expected async = false")
	}

	post := byEvent["PostToolUse"]
	if !post.Async {
		t.Error("expected async = true")
	}
	if post.Timeout != 10 {
		t.Errorf("timeout = %d, want 10", post.Timeout)
	}
}

func TestExtractHooks_Empty(t *testing.T) {
	settings := &claude.Settings{}
	hooks, err := claude.ExtractHooks(settings, "global")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hooks) != 0 {
		t.Errorf("expected empty slice, got %v", hooks)
	}
}
