package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
)

// TestUpdateSearch_CharInput tests the default branch of updateSearch
// which calls applySearch when searching is true.
func TestUpdateSearch_CharInputFilters(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
		makeChannel("3", "gamma"),
	}
	m := newModel(channels, defaultCfg())
	m.searching = true

	// Typing a character goes through default branch → applySearch
	m2, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	result := m2.(model)

	if !result.searching {
		t.Error("should remain in search mode after typing a character")
	}
}

func TestUpdateSearch_EscRestoresCursor(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
		makeChannel("3", "gamma"),
	}
	m := newModel(channels, defaultCfg())
	m.searching = true
	m.prevCursor = 2
	m.channelCursor = 0

	m2, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyEscape})
	result := m2.(model)

	if result.channelCursor != 2 {
		t.Errorf("after esc: cursor = %d, want 2 (prevCursor)", result.channelCursor)
	}
}

func TestUpdateSearch_EscRestoresFullFilter(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.searching = true
	m.filtered = []int{0} // simulate filtered down

	m2, _ := m.updateSearch(tea.KeyMsg{Type: tea.KeyEscape})
	result := m2.(model)

	if len(result.filtered) != 2 {
		t.Errorf("after esc: filtered len = %d, want 2 (all channels)", len(result.filtered))
	}
}

func TestApplySearch_EmptyQueryResetsFilter(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.filtered = []int{0} // simulate filtered

	m.searchInput.SetValue("")
	m.applySearch()

	if len(m.filtered) != 2 {
		t.Errorf("applySearch with empty query: filtered len = %d, want 2", len(m.filtered))
	}
}

func TestApplySearch_MatchingQuery(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("alpha-id", "alpha"),
		makeChannel("beta-id", "beta"),
		makeChannel("gamma-id", "gamma"),
	}
	m := newModel(channels, defaultCfg())

	m.searchInput.SetValue("alpha")
	m.applySearch()

	if len(m.filtered) != 1 {
		t.Errorf("applySearch('alpha'): filtered len = %d, want 1", len(m.filtered))
	}
}

func TestApplySearch_NoMatch(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())

	m.searchInput.SetValue("zzzzzzz")
	m.applySearch()

	if len(m.filtered) != 0 {
		t.Errorf("applySearch('zzzzzzz'): filtered len = %d, want 0", len(m.filtered))
	}
	if m.channelCursor != 0 {
		t.Errorf("applySearch no match: cursor = %d, want 0", m.channelCursor)
	}
}

func TestApplyContentSearch_EmptyQuery(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.contentSearching = true
	m.width = 80
	m.height = 24

	m.searchInput.SetValue("")
	m.applyContentSearch()

	if len(m.contentMatches) != 0 {
		t.Errorf("applyContentSearch with empty query: matches = %d, want 0", len(m.contentMatches))
	}
}

func TestRenderContentWithHighlight_NoMatches(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.contentMatches = nil

	// Should not panic
	m.renderContentWithHighlight()
}
