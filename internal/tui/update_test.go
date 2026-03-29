package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
)

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func specialKey(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func TestUpdateNormal_CursorDown(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
		makeChannel("3", "gamma"),
	}
	m := newModel(channels, defaultCfg())
	m.focus = listPanel

	m2, _ := m.updateNormal(runeKey('j'))
	result := m2.(model)
	if result.channelCursor != 1 {
		t.Errorf("after j: cursor = %d, want 1", result.channelCursor)
	}
}

func TestUpdateNormal_CursorDownArrow(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.focus = listPanel

	m2, _ := m.updateNormal(specialKey(tea.KeyDown))
	result := m2.(model)
	if result.channelCursor != 1 {
		t.Errorf("after down: cursor = %d, want 1", result.channelCursor)
	}
}

func TestUpdateNormal_CursorDownAtEnd(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.focus = listPanel
	m.channelCursor = 1 // already at end

	m2, _ := m.updateNormal(runeKey('j'))
	result := m2.(model)
	if result.channelCursor != 1 {
		t.Errorf("after j at end: cursor = %d, want 1", result.channelCursor)
	}
}

func TestUpdateNormal_CursorUp(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.focus = listPanel
	m.channelCursor = 1

	m2, _ := m.updateNormal(runeKey('k'))
	result := m2.(model)
	if result.channelCursor != 0 {
		t.Errorf("after k: cursor = %d, want 0", result.channelCursor)
	}
}

func TestUpdateNormal_CursorUpAtTop(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = listPanel
	m.channelCursor = 0

	m2, _ := m.updateNormal(runeKey('k'))
	result := m2.(model)
	if result.channelCursor != 0 {
		t.Errorf("after k at top: cursor = %d, want 0", result.channelCursor)
	}
}

func TestUpdateNormal_CursorUpArrow(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.focus = listPanel
	m.channelCursor = 1

	m2, _ := m.updateNormal(specialKey(tea.KeyUp))
	result := m2.(model)
	if result.channelCursor != 0 {
		t.Errorf("after up: cursor = %d, want 0", result.channelCursor)
	}
}

func TestUpdateNormal_TabTogglesFocus(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = listPanel

	m2, _ := m.updateNormal(specialKey(tea.KeyTab))
	result := m2.(model)
	if result.focus != detailPanel {
		t.Errorf("after Tab from list: focus = %v, want detailPanel", result.focus)
	}

	m3, _ := result.updateNormal(specialKey(tea.KeyTab))
	result = m3.(model)
	if result.focus != listPanel {
		t.Errorf("after Tab from detail: focus = %v, want listPanel", result.focus)
	}
}

func TestUpdateNormal_ShiftTabTogglesFocus(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel

	m2, _ := m.updateNormal(specialKey(tea.KeyShiftTab))
	result := m2.(model)
	if result.focus != listPanel {
		t.Errorf("after ShiftTab from detail: focus = %v, want listPanel", result.focus)
	}
}

func TestUpdateNormal_RightFromListFocusesDetail(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = listPanel

	m2, _ := m.updateNormal(specialKey(tea.KeyRight))
	result := m2.(model)
	if result.focus != detailPanel {
		t.Errorf("after right from list: focus = %v, want detailPanel", result.focus)
	}
}

func TestUpdateNormal_RightCyclesDetailTabs(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel
	m.detailTab = TabSettings

	m2, _ := m.updateNormal(specialKey(tea.KeyRight))
	result := m2.(model)
	if result.detailTab != TabClaudeMD {
		t.Errorf("after right in detail: tab = %v, want TabClaudeMD(%d)", result.detailTab, TabClaudeMD)
	}
}

func TestUpdateNormal_RightTabWrapsAround(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel
	m.detailTab = DetailTab(len(detailTabNames) - 1) // last tab

	m2, _ := m.updateNormal(specialKey(tea.KeyRight))
	result := m2.(model)
	if result.detailTab != TabSettings {
		t.Errorf("after right wrap: tab = %v, want TabSettings(0)", result.detailTab)
	}
}

func TestUpdateNormal_LeftCyclesDetailTabs(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel
	m.detailTab = TabClaudeMD

	m2, _ := m.updateNormal(specialKey(tea.KeyLeft))
	result := m2.(model)
	if result.detailTab != TabSettings {
		t.Errorf("after left in detail: tab = %v, want TabSettings(0)", result.detailTab)
	}
}

func TestUpdateNormal_LeftTabWrapsAround(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel
	m.detailTab = TabSettings // first tab

	m2, _ := m.updateNormal(specialKey(tea.KeyLeft))
	result := m2.(model)
	want := DetailTab(len(detailTabNames) - 1)
	if result.detailTab != want {
		t.Errorf("after left wrap: tab = %v, want %v", result.detailTab, want)
	}
}

func TestUpdateNormal_SlashEntersSearch(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())

	m2, _ := m.updateNormal(runeKey('/'))
	result := m2.(model)
	if !result.searching {
		t.Error("after /: searching should be true")
	}
	if result.contentSearching {
		t.Error("after /: contentSearching should be false")
	}
}

func TestUpdateNormal_QuestionMarkEntersContentSearch(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())

	m2, _ := m.updateNormal(runeKey('?'))
	result := m2.(model)
	if !result.contentSearching {
		t.Error("after ?: contentSearching should be true")
	}
	if result.searching {
		t.Error("after ?: searching should be false")
	}
}

func TestUpdateNormal_CmdEnterNavigatesAndQuits(t *testing.T) {
	ch := makeChannel("1", "myproject")
	ch.Path = "/home/user/myproject"
	m := newModel([]channel.Channel{ch}, defaultCfg())

	m2, cmd := m.updateNormal(tea.KeyMsg{Type: tea.KeyEnter, Alt: true})
	result := m2.(model)
	if result.navigateTo != "/home/user/myproject" {
		t.Errorf("navigateTo = %q, want /home/user/myproject", result.navigateTo)
	}
	if cmd == nil {
		t.Error("cmd should be non-nil (tea.Quit) after alt+enter")
	}
}

func TestUpdateNormal_CmdEnterNoopWithEmptyList(t *testing.T) {
	m := newModel([]channel.Channel{}, defaultCfg())

	m2, cmd := m.updateNormal(tea.KeyMsg{Type: tea.KeyEnter, Alt: true})
	result := m2.(model)
	if result.navigateTo != "" {
		t.Errorf("navigateTo should be empty, got %q", result.navigateTo)
	}
	if cmd != nil {
		t.Error("cmd should be nil when no channel selected")
	}
}

func TestUpdateNormal_QuitReturnsQuitCmd(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())

	_, cmd := m.updateNormal(runeKey('q'))
	if cmd == nil {
		t.Error("q should return quit cmd")
	}
}

func TestUpdateSearch_EscCancelsSearch(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.searching = true
	m.prevCursor = 1

	m2, _ := m.updateSearch(specialKey(tea.KeyEscape))
	result := m2.(model)

	if result.searching {
		t.Error("after esc: searching should be false")
	}
	if result.channelCursor != 1 {
		t.Errorf("after esc: cursor = %d, want 1 (prevCursor)", result.channelCursor)
	}
	if len(result.filtered) != 2 {
		t.Errorf("after esc: filtered len = %d, want 2", len(result.filtered))
	}
}

func TestUpdateSearch_EscCancelsContentSearch(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.contentSearching = true

	m2, _ := m.updateSearch(specialKey(tea.KeyEscape))
	result := m2.(model)

	if result.contentSearching {
		t.Error("after esc: contentSearching should be false")
	}
}

func TestUpdateSearch_EnterConfirmsSearch(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.searching = true

	m2, _ := m.updateSearch(specialKey(tea.KeyEnter))
	result := m2.(model)

	if result.searching {
		t.Error("after enter: searching should be false")
	}
}

func TestEditTargetForTab(t *testing.T) {
	ch := &channel.Channel{Path: "/project"}

	tests := []struct {
		tab  DetailTab
		want string
	}{
		{TabClaudeMD, "/project/CLAUDE.md"},
		{TabSettings, "/project/.claude/settings.json"},
		{TabHooks, "/project/.claude/settings.json"},
		{TabMCP, "/project/.claude/settings.json"},
		{TabHealth, "/project/.claude/settings.json"},
	}

	for _, tt := range tests {
		got := editTargetForTab(ch, tt.tab)
		if got != tt.want {
			t.Errorf("editTargetForTab(tab=%v) = %q, want %q", tt.tab, got, tt.want)
		}
	}
}

func TestListWidth(t *testing.T) {
	tests := []struct {
		windowWidth int
		wantMin     int
		wantMax     int
	}{
		{0, 28, 28},    // below minimum → clamped to 28
		{100, 28, 28},  // 25% of 100 = 25, clamped to 28
		{200, 36, 36},  // 25% of 200 = 50, clamped to 36
	}

	for _, tt := range tests {
		m := newModel([]channel.Channel{}, defaultCfg())
		m.width = tt.windowWidth
		got := m.listWidth()
		if got < tt.wantMin || got > tt.wantMax {
			t.Errorf("listWidth(width=%d) = %d, want [%d, %d]", tt.windowWidth, got, tt.wantMin, tt.wantMax)
		}
	}
}
