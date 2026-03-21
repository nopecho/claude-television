package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderSettingsTab(ch *channel.Channel) string {
	var b strings.Builder

	b.WriteString(section("Project"))
	b.WriteString(kv("path", ch.Path) + "\n")
	b.WriteString(kv("status", string(ch.Status)) + "\n")

	if ch.Data.Settings != nil {
		b.WriteString(renderSettingsSection(ch.Data.Settings, "Project Settings"))
	}
	if ch.Data.LocalSettings != nil {
		b.WriteString(renderSettingsSection(ch.Data.LocalSettings, "Local Settings (override)"))
	}
	if ch.Data.Settings == nil && ch.Data.LocalSettings == nil {
		b.WriteString(section("Settings"))
		b.WriteString("    No settings.json found\n")
	}

	if !ch.IsGlobal {
		globalCh := m.findGlobalChannel()
		if globalCh != nil && globalCh.Data != nil && globalCh.Data.Settings != nil && ch.Data.Settings != nil {
			diffs := diffSettings(globalCh.Data.Settings, ch.Data.Settings)
			if len(diffs) > 0 {
				b.WriteString(section("Overrides from Global"))
				for _, d := range diffs {
					b.WriteString(fmt.Sprintf("    %s: %s → %s\n",
						labelStyle.Render(d.key),
						labelStyle.Render(d.globalVal),
						valueStyle.Render(d.projectVal)))
				}
			}
		}
	}

	return b.String()
}

type settingsDiff struct {
	key        string
	globalVal  string
	projectVal string
}

func diffSettings(global, project *claude.Settings) []settingsDiff {
	var diffs []settingsDiff
	if project.Model != "" && project.Model != global.Model {
		diffs = append(diffs, settingsDiff{"model", global.Model, project.Model})
	}
	if project.Language != "" && project.Language != global.Language {
		diffs = append(diffs, settingsDiff{"language", global.Language, project.Language})
	}
	if project.TeammateMode != "" && project.TeammateMode != global.TeammateMode {
		diffs = append(diffs, settingsDiff{"teammateMode", global.TeammateMode, project.TeammateMode})
	}
	return diffs
}

func renderSettingsSection(s *claude.Settings, title string) string {
	var b strings.Builder
	b.WriteString(section(title))
	if s.Model != "" {
		b.WriteString(kv("model", s.Model) + "\n")
	}
	if s.Language != "" {
		b.WriteString(kv("language", s.Language) + "\n")
	}
	if s.TeammateMode != "" {
		b.WriteString(kv("teammate", s.TeammateMode) + "\n")
	}
	if s.PlansDirectory != "" {
		b.WriteString(kv("plans dir", s.PlansDirectory) + "\n")
	}
	if len(s.Env) > 0 {
		b.WriteString(section("Environment"))
		for k, v := range s.Env {
			b.WriteString(fmt.Sprintf("    %s = %s\n", k, v))
		}
	}
	if len(s.Permissions.Allow) > 0 || len(s.Permissions.Deny) > 0 {
		b.WriteString(section("Permissions"))
		if len(s.Permissions.Allow) > 0 {
			b.WriteString(fmt.Sprintf("    Allow (%d):\n", len(s.Permissions.Allow)))
			for _, p := range s.Permissions.Allow {
				b.WriteString(fmt.Sprintf("      %s %s\n", statusHealthy, p))
			}
		}
		if len(s.Permissions.Deny) > 0 {
			b.WriteString(fmt.Sprintf("    Deny (%d):\n", len(s.Permissions.Deny)))
			for _, p := range s.Permissions.Deny {
				b.WriteString(fmt.Sprintf("      %s %s\n", statusError, p))
			}
		}
	}
	return b.String()
}
