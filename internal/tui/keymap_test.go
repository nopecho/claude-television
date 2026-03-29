package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestParseKey(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyMsg
		want keyAction
	}{
		{"q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, keyQuit},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}, keyQuit},
		{"up", tea.KeyMsg{Type: tea.KeyUp}, keyUp},
		{"k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}, keyUp},
		{"down", tea.KeyMsg{Type: tea.KeyDown}, keyDown},
		{"j", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, keyDown},
		{"left", tea.KeyMsg{Type: tea.KeyLeft}, keyLeft},
		{"h", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, keyLeft},
		{"right", tea.KeyMsg{Type: tea.KeyRight}, keyRight},
		{"l", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, keyRight},
		{"tab", tea.KeyMsg{Type: tea.KeyTab}, keyTab},
		{"shift+tab", tea.KeyMsg{Type: tea.KeyShiftTab}, keyShiftTab},
		{"enter", tea.KeyMsg{Type: tea.KeyEnter}, keyEnter},
		{"alt+enter", tea.KeyMsg{Type: tea.KeyEnter, Alt: true}, keyCmdEnter},
		{"/", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}, keySlash},
		{"esc", tea.KeyMsg{Type: tea.KeyEscape}, keyEscape},
		{"p", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}, keyPin},
		{"e", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}, keyEdit},
		{"ctrl+d", tea.KeyMsg{Type: tea.KeyCtrlD}, keyScrollDown},
		{"ctrl+u", tea.KeyMsg{Type: tea.KeyCtrlU}, keyScrollUp},
		{"?", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, keyContentSearch},
		{"g", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}, keyGroup},
		{"unknown_x", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}, keyNone},
		{"unknown_z", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}}, keyNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseKey(tt.msg)
			if got != tt.want {
				t.Errorf("parseKey(%q) = %v, want %v", tt.msg.String(), got, tt.want)
			}
		})
	}
}
