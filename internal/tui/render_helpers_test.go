package tui

import (
	"strings"
	"testing"

	"github.com/nopecho/claude-television/internal/channel"
)

func TestBoolIcon(t *testing.T) {
	trueIcon := boolIcon(true)
	falseIcon := boolIcon(false)

	if trueIcon == "" {
		t.Error("boolIcon(true) should not be empty")
	}
	if falseIcon == "" {
		t.Error("boolIcon(false) should not be empty")
	}
	if trueIcon == falseIcon {
		t.Error("boolIcon(true) and boolIcon(false) should differ")
	}
}

func TestSectionEmpty(t *testing.T) {
	got := sectionEmpty()
	if got == "" {
		t.Error("sectionEmpty() should not be empty")
	}
}

func TestBullet(t *testing.T) {
	got := bullet("my item")
	if !strings.Contains(got, "my item") {
		t.Errorf("bullet() should contain the item text, got: %q", got)
	}
	if !strings.Contains(got, "•") {
		t.Errorf("bullet() should contain bullet symbol, got: %q", got)
	}
}

func TestHelpEntry(t *testing.T) {
	got := helpEntry("j/k", "move")
	if !strings.Contains(got, "j/k") {
		t.Errorf("helpEntry() should contain key, got: %q", got)
	}
	if !strings.Contains(got, "move") {
		t.Errorf("helpEntry() should contain description, got: %q", got)
	}
}

func TestOrderedGroup_PreservesInsertionOrder(t *testing.T) {
	items := []string{"b", "a", "b", "c", "a"}
	order, groups := orderedGroup(items, func(s string) string { return s })

	if len(order) != 3 {
		t.Fatalf("order len = %d, want 3", len(order))
	}
	if order[0] != "b" || order[1] != "a" || order[2] != "c" {
		t.Errorf("order = %v, want [b a c]", order)
	}
	if len(groups["b"]) != 2 {
		t.Errorf("groups[b] len = %d, want 2", len(groups["b"]))
	}
	if len(groups["a"]) != 2 {
		t.Errorf("groups[a] len = %d, want 2", len(groups["a"]))
	}
	if len(groups["c"]) != 1 {
		t.Errorf("groups[c] len = %d, want 1", len(groups["c"]))
	}
}

func TestOrderedGroup_Empty(t *testing.T) {
	order, groups := orderedGroup([]string{}, func(s string) string { return s })
	if len(order) != 0 {
		t.Errorf("order len = %d, want 0", len(order))
	}
	if len(groups) != 0 {
		t.Errorf("groups len = %d, want 0", len(groups))
	}
}

func TestStatusIconStr(t *testing.T) {
	tests := []struct {
		status channel.ChannelStatus
	}{
		{channel.StatusHealthy},
		{channel.StatusWarning},
		{channel.StatusError},
		{"unknown"},
	}

	for _, tt := range tests {
		got := statusIconStr(tt.status)
		if got == "" {
			t.Errorf("statusIconStr(%q) should not be empty", tt.status)
		}
	}
}

func TestStatusIconStr_Distinct(t *testing.T) {
	healthy := statusIconStr(channel.StatusHealthy)
	warning := statusIconStr(channel.StatusWarning)
	errIcon := statusIconStr(channel.StatusError)

	if healthy == warning {
		t.Error("healthy and warning icons should differ")
	}
	if healthy == errIcon {
		t.Error("healthy and error icons should differ")
	}
	if warning == errIcon {
		t.Error("warning and error icons should differ")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s    string
		max  int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hell…"},
		{"ab", 2, "ab"},
		{"abc", 2, "a…"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.s, tt.max)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.s, tt.max, got, tt.want)
		}
	}
}

func TestSectionHeader(t *testing.T) {
	got := sectionHeader("My Section")
	if !strings.Contains(got, "My Section") {
		t.Errorf("sectionHeader() should contain title, got: %q", got)
	}
	if !strings.Contains(got, "▎") {
		t.Errorf("sectionHeader() should contain bar symbol, got: %q", got)
	}
}

func TestSectionLine(t *testing.T) {
	got := sectionLine("content here")
	if !strings.Contains(got, "content here") {
		t.Errorf("sectionLine() should contain content, got: %q", got)
	}
	if !strings.Contains(got, "│") {
		t.Errorf("sectionLine() should contain bar symbol, got: %q", got)
	}
}

func TestSection(t *testing.T) {
	got := section("Title")
	if !strings.Contains(got, "Title") {
		t.Errorf("section() should contain title, got: %q", got)
	}
}

func TestKv(t *testing.T) {
	got := kv("key", "value", 10)
	if !strings.Contains(got, "key") {
		t.Errorf("kv() should contain key, got: %q", got)
	}
	if !strings.Contains(got, "value") {
		t.Errorf("kv() should contain value, got: %q", got)
	}
}

func TestEmptyState(t *testing.T) {
	got := emptyState("Title", "message", "hint")
	if !strings.Contains(got, "message") {
		t.Errorf("emptyState() should contain message, got: %q", got)
	}
	if !strings.Contains(got, "hint") {
		t.Errorf("emptyState() should contain hint, got: %q", got)
	}
}

func TestEmptyState_NoHint(t *testing.T) {
	got := emptyState("Title", "message", "")
	if !strings.Contains(got, "message") {
		t.Errorf("emptyState() should contain message, got: %q", got)
	}
}
