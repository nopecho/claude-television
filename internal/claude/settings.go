package claude

import (
	"encoding/json"
	"fmt"
	"os"
)

type Settings struct {
	Model          string                `json:"model"`
	Language       string                `json:"language"`
	Env            map[string]string     `json:"env"`
	Permissions    Permissions           `json:"permissions"`
	EnabledPlugins map[string]bool       `json:"enabledPlugins"`
	Hooks          map[string][]HookRule `json:"hooks"`
	StatusLine     json.RawMessage       `json:"statusLine"`
	PlansDirectory string                `json:"plansDirectory"`
	TeammateMode   string                `json:"teammateMode"`
	Raw            map[string]any        `json:"-"`
}

type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

type HookRule struct {
	Matcher string       `json:"matcher"`
	Hooks   []HookAction `json:"hooks"`
}

type HookAction struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Prompt  string `json:"prompt"`
	Async   bool   `json:"async"`
	Timeout int    `json:"timeout"`
}

func ParseSettings(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse settings: %w", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse raw settings: %w", err)
	}
	s.Raw = raw
	return &s, nil
}
