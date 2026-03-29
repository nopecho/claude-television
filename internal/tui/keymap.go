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
	keyEdit
	keyScrollDown
	keyScrollUp
	keyContentSearch
	keyGroup
	keyTab1
	keyTab2
	keyTab3
	keyTab4
	keyTab5
	keyTab6
	keyTab7
	keyTab8
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
	case "e":
		return keyEdit
	case "ctrl+d":
		return keyScrollDown
	case "ctrl+u":
		return keyScrollUp
	case "?":
		return keyContentSearch
	case "g":
		return keyGroup
	case "1":
		return keyTab1
	case "2":
		return keyTab2
	case "3":
		return keyTab3
	case "4":
		return keyTab4
	case "5":
		return keyTab5
	case "6":
		return keyTab6
	case "7":
		return keyTab7
	case "8":
		return keyTab8
	}
	return keyNone
}
