package tui

import (
	"strings"
	"testing"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func TestView_Loading(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	// width=0 → loading state

	got := m.View()
	if !strings.Contains(got, "Loading") {
		t.Errorf("View() with width=0 should show loading, got: %q", got)
	}
}

func TestView_WithSize(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.width = 100
	m.height = 30

	got := m.View()
	if !strings.Contains(got, "ctv") {
		t.Errorf("View() should contain app title 'ctv', got: %q", got)
	}
	if !strings.Contains(got, "alpha") {
		t.Errorf("View() should contain channel name, got: %q", got)
	}
}

func TestView_SearchMode(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.width = 100
	m.height = 30
	m.searching = true

	got := m.View()
	if !strings.Contains(got, "matches") {
		t.Errorf("View() in search mode should show match count, got: %q", got)
	}
}

func TestView_ContentSearchMode(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.width = 100
	m.height = 30
	m.contentSearching = true

	got := m.View()
	if !strings.Contains(got, "content search") {
		t.Errorf("View() in content search mode should show hint, got: %q", got)
	}
}

func TestView_GroupingMode(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.width = 100
	m.height = 30
	m.grouping = true

	got := m.View()
	if !strings.Contains(got, "group") {
		t.Errorf("View() in grouping mode should show group hint, got: %q", got)
	}
}

func TestView_DetailPanelFocused(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.width = 100
	m.height = 30
	m.focus = detailPanel

	got := m.View()
	// Detail panel focus should render without panicking
	if got == "" {
		t.Error("View() should not return empty string")
	}
}

func TestCountStatuses_Counts(t *testing.T) {
	channels := []channel.Channel{
		{ID: "1", Name: "h1", Status: channel.StatusHealthy, Data: &channel.ChannelData{}},
		{ID: "2", Name: "h2", Status: channel.StatusHealthy, Data: &channel.ChannelData{}},
		{ID: "3", Name: "e1", Status: channel.StatusError, Data: &channel.ChannelData{}},
	}
	m := newModel(channels, defaultCfg())

	healthy, _, errCount := m.countStatuses()
	if healthy != 2 {
		t.Errorf("countStatuses() healthy should be 2, got: %d", healthy)
	}
	if errCount != 1 {
		t.Errorf("countStatuses() errCount should be 1, got: %d", errCount)
	}
}

func TestCountStatuses_GlobalSkipped(t *testing.T) {
	channels := []channel.Channel{
		{ID: "1", Name: "proj", Status: channel.StatusHealthy, Data: &channel.ChannelData{}},
		{ID: "2", Name: "global", Status: channel.StatusHealthy, IsGlobal: true, Data: &channel.ChannelData{}},
	}
	m := newModel(channels, defaultCfg())

	healthy, _, _ := m.countStatuses()
	if healthy != 1 {
		t.Errorf("countStatuses() should skip global channel, healthy=%d", healthy)
	}
}

func TestCountStatuses_WithWarning(t *testing.T) {
	channels := []channel.Channel{
		{
			ID: "1", Name: "proj", Status: channel.StatusWarning,
			Data: &channel.ChannelData{
				HealthIssues: []claude.HealthIssue{{Message: "missing settings"}},
			},
		},
	}
	m := newModel(channels, defaultCfg())

	_, warning, _ := m.countStatuses()
	if warning != 1 {
		t.Errorf("countStatuses() warning should be 1, got: %d", warning)
	}
}

func TestRenderDetailTabs_ContainsAllTabs(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())

	got := m.renderDetailTabs()
	for _, name := range detailTabNames {
		if !strings.Contains(got, name) {
			t.Errorf("renderDetailTabs() should contain tab %q, got: %q", name, got)
		}
	}
}

func TestRenderDetailTabs_ActiveTabHighlighted(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.detailTab = TabSettings

	got := m.renderDetailTabs()
	// Active tab should contain tab number prefix
	if !strings.Contains(got, "1") {
		t.Errorf("renderDetailTabs() should show tab numbers, got: %q", got)
	}
}

func TestRenderHelpBar_ListPanel(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = listPanel

	got := m.renderHelpBar()
	if !strings.Contains(got, "quit") {
		t.Errorf("help bar should contain 'quit', got: %q", got)
	}
	if !strings.Contains(got, "search") {
		t.Errorf("help bar should contain 'search', got: %q", got)
	}
}

func TestRenderHelpBar_DetailPanel(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.focus = detailPanel

	got := m.renderHelpBar()
	if !strings.Contains(got, "tabs") {
		t.Errorf("detail help bar should contain 'tabs', got: %q", got)
	}
	if !strings.Contains(got, "scroll") {
		t.Errorf("detail help bar should contain 'scroll', got: %q", got)
	}
}

func TestRenderHelpBar_SearchMode(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.searching = true

	got := m.renderHelpBar()
	if !strings.Contains(got, "cancel") {
		t.Errorf("search help bar should contain 'cancel', got: %q", got)
	}
	if !strings.Contains(got, "confirm") {
		t.Errorf("search help bar should contain 'confirm', got: %q", got)
	}
}

func TestInjectBorderTitle_EmptyTitle(t *testing.T) {
	box := "some box content"
	got := injectBorderTitle(box, "")
	if got != box {
		t.Errorf("injectBorderTitle with empty title should return unchanged box")
	}
}

func TestInjectBorderFooter_EmptyFooter(t *testing.T) {
	box := "some\nbox\ncontent"
	got := injectBorderFooter(box, "")
	if got != box {
		t.Errorf("injectBorderFooter with empty footer should return unchanged box")
	}
}

func TestInjectBorderFooter_ShortBox(t *testing.T) {
	box := "single"
	got := injectBorderFooter(box, "footer")
	if got != box {
		t.Errorf("injectBorderFooter with single-line box should return unchanged box")
	}
}
