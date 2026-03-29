package tui

import (
	"strings"
	"testing"

	"github.com/nopecho/claude-television/internal/channel"
)

func TestRenderChannelList_Empty(t *testing.T) {
	m := newModel([]channel.Channel{}, defaultCfg())

	got := m.renderChannelList(10)
	if !strings.Contains(got, "No channels found") {
		t.Errorf("empty list should show 'No channels found', got: %q", got)
	}
}

func TestRenderChannelList_ShowsChannelNames(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(10)
	if !strings.Contains(got, "alpha") {
		t.Errorf("channel list should contain 'alpha', got: %q", got)
	}
	if !strings.Contains(got, "beta") {
		t.Errorf("channel list should contain 'beta', got: %q", got)
	}
}

func TestRenderChannelList_SelectedIndicator(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "alpha"),
		makeChannel("2", "beta"),
	}
	m := newModel(channels, defaultCfg())
	m.channelCursor = 0

	got := m.renderChannelList(10)
	if !strings.Contains(got, "▸") {
		t.Errorf("selected channel should have ▸ indicator, got: %q", got)
	}
}

func TestRenderChannelList_GroupHeader(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "project-a", withGroup("work")),
		makeChannel("2", "project-b", withGroup("work")),
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(10)
	if !strings.Contains(got, "work") {
		t.Errorf("grouped channel list should show group header, got: %q", got)
	}
}

func TestRenderChannelList_GlobalChannelIcon(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "global-channel", withGlobal()),
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(10)
	if !strings.Contains(got, "⚙") {
		t.Errorf("global channel should show ⚙ icon, got: %q", got)
	}
}

func TestRenderChannelList_HeightLimitsOutput(t *testing.T) {
	channels := make([]channel.Channel, 20)
	for i := range channels {
		channels[i] = makeChannel(string(rune('a'+i)), string(rune('a'+i))+"-project")
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(3)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) > 3 {
		t.Errorf("renderChannelList(height=3) returned %d lines, want <= 3", len(lines))
	}
}

func TestRenderChannelList_TruncatesLongNames(t *testing.T) {
	longName := "this-is-a-very-long-project-name-that-exceeds-limit"
	channels := []channel.Channel{
		makeChannel("1", longName),
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(10)
	// truncated name should not appear in full
	if strings.Contains(got, longName) {
		t.Errorf("long channel name should be truncated, but got full name in: %q", got)
	}
	// should contain ellipsis
	if !strings.Contains(got, "…") {
		t.Errorf("truncated name should contain ellipsis, got: %q", got)
	}
}

func TestRenderChannelList_SingleChannel(t *testing.T) {
	channels := []channel.Channel{
		makeChannel("1", "solo"),
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(5)
	if !strings.Contains(got, "solo") {
		t.Errorf("should contain 'solo', got: %q", got)
	}
}

func TestRenderChannelList_StatusHealthy(t *testing.T) {
	channels := []channel.Channel{
		{ID: "1", Name: "healthy-ch", Path: "/tmp", Status: channel.StatusHealthy, Data: &channel.ChannelData{}},
	}
	m := newModel(channels, defaultCfg())

	got := m.renderChannelList(5)
	if !strings.Contains(got, "healthy-ch") {
		t.Errorf("should contain channel name, got: %q", got)
	}
}
