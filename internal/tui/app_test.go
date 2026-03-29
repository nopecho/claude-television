package tui

import (
	"testing"

	"github.com/nopecho/claude-television/internal/channel"
)

func TestNewModel_InitialState(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())

	if m.channelCursor != 0 {
		t.Errorf("initial cursor = %d, want 0", m.channelCursor)
	}
	if len(m.filtered) != 2 {
		t.Errorf("filtered len = %d, want 2", len(m.filtered))
	}
	if m.focus != listPanel {
		t.Errorf("initial focus = %v, want listPanel", m.focus)
	}
	if m.searching {
		t.Error("initial searching should be false")
	}
	if m.grouping {
		t.Error("initial grouping should be false")
	}
}

func TestSortChannels_PinnedFirst(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "ungrouped"),
		makeChannel("2", "pinned", withPinned(true)),
		makeChannel("3", "grouped", withGroup("work")),
	}
	m := newModel(channels, defaultCfg())

	if !m.channels[0].Pinned {
		t.Errorf("first channel should be pinned, got %q", m.channels[0].Name)
	}
	if m.channels[1].Group != "work" {
		t.Errorf("second channel should be grouped, got %q", m.channels[1].Name)
	}
	if m.channels[2].Name != "ungrouped" {
		t.Errorf("third channel should be ungrouped, got %q", m.channels[2].Name)
	}
}

func TestSortChannels_GroupOrder(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "w1", withGroup("work")),
		makeChannel("2", "h1", withGroup("home")),
		makeChannel("3", "h2", withGroup("home")),
	}
	m := newModel(channels, defaultCfg())

	// work group appears before home (first-seen order)
	if m.channels[0].Group != "work" {
		t.Errorf("first group should be work, got %q", m.channels[0].Group)
	}
	if m.channels[1].Group != "home" {
		t.Errorf("second group should be home, got %q", m.channels[1].Group)
	}
	if m.channels[2].Group != "home" {
		t.Errorf("third group should be home, got %q", m.channels[2].Group)
	}
}

func TestSortChannels_AllCategories(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "bare"),
		makeChannel("2", "grouped", withGroup("g")),
		makeChannel("3", "starred", withPinned(true)),
	}
	m := newModel(channels, defaultCfg())

	// order: pinned → grouped → ungrouped
	if !m.channels[0].Pinned {
		t.Errorf("position 0 should be pinned, got name=%q", m.channels[0].Name)
	}
	if m.channels[1].Group != "g" {
		t.Errorf("position 1 should be grouped, got name=%q group=%q", m.channels[1].Name, m.channels[1].Group)
	}
	if m.channels[2].Name != "bare" {
		t.Errorf("position 2 should be bare, got name=%q", m.channels[2].Name)
	}
}

func TestSelectedChannel(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())

	ch := m.selectedChannel()
	if ch == nil {
		t.Fatal("selectedChannel() returned nil at cursor 0")
	}
	if ch.Name != "alpha" {
		t.Errorf("selectedChannel() = %q, want alpha", ch.Name)
	}

	m.channelCursor = 1
	ch = m.selectedChannel()
	if ch == nil {
		t.Fatal("selectedChannel() returned nil at cursor 1")
	}
	if ch.Name != "beta" {
		t.Errorf("selectedChannel() = %q, want beta", ch.Name)
	}
}

func TestSelectedChannel_Empty(t *testing.T) {
	m := newModel([]channel.Channel{}, defaultCfg())
	if m.selectedChannel() != nil {
		t.Error("selectedChannel() should return nil for empty channel list")
	}
}

func TestSelectedChannel_OutOfBounds(t *testing.T) {
	m := newModel([]channel.Channel{makeChannel("1", "alpha")}, defaultCfg())
	m.channelCursor = 99
	if m.selectedChannel() != nil {
		t.Error("selectedChannel() should return nil when cursor is out of bounds")
	}
}

func TestFindGlobalChannel(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "project"),
		makeChannel("2", "global", withGlobal()),
	}
	m := newModel(channels, defaultCfg())

	ch := m.findGlobalChannel()
	if ch == nil {
		t.Fatal("findGlobalChannel() returned nil")
	}
	if !ch.IsGlobal {
		t.Error("findGlobalChannel() returned non-global channel")
	}
}

func TestFindGlobalChannel_None(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "project"),
	}
	m := newModel(channels, defaultCfg())

	if m.findGlobalChannel() != nil {
		t.Error("findGlobalChannel() should return nil when no global channel exists")
	}
}

func TestResetFilter(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
		makeChannel("3", "gamma"),
	}
	m := newModel(channels, defaultCfg())
	m.filtered = []int{0} // simulate filtered state

	m.resetFilter()

	if len(m.filtered) != 3 {
		t.Errorf("resetFilter() len = %d, want 3", len(m.filtered))
	}
	for i, idx := range m.filtered {
		if idx != i {
			t.Errorf("filtered[%d] = %d, want %d", i, idx, i)
		}
	}
}
