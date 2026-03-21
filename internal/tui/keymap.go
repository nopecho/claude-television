package tui

import tea "github.com/charmbracelet/bubbletea"

type keyAction int

const (
	keyNone keyAction = iota
	keyQuit
	keyUp
	keyDown
	keyLeft
	keyRight
	keyTab
	keyShiftTab
	keyEnter
	keyCmdEnter
	keySlash
	keyEscape
	keyPin
	keyGroup
	keyEdit
)

func parseKey(msg tea.KeyMsg) keyAction {
	switch msg.String() {
	case "q", "ctrl+c":
		return keyQuit
	case "up", "k":
		return keyUp
	case "down", "j":
		return keyDown
	case "left", "h":
		return keyLeft
	case "right", "l":
		return keyRight
	case "tab":
		return keyTab
	case "shift+tab":
		return keyShiftTab
	case "enter":
		return keyEnter
	case "alt+enter":
		return keyCmdEnter
	case "/":
		return keySlash
	case "esc":
		return keyEscape
	case "p":
		return keyPin
	case "g":
		return keyGroup
	case "e":
		return keyEdit
	}
	return keyNone
}
