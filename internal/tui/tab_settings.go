package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderSettingsTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	// Project card
	b.WriteString(card("Project", []string{
		cardKV("path", ch.Path, 8),
		cardKV("status", statusIconStr(ch.Status)+" "+string(ch.Status), 8),
	}, w))
	b.WriteString("\n")

	if ch.Data != nil && ch.Data.Settings != nil {
		b.WriteString(renderSettingsCard(ch.Data.Settings, "Project Settings", w))
		b.WriteString("\n")
	}
	if ch.Data != nil && ch.Data.LocalSettings != nil {
		b.WriteString(renderSettingsCard(ch.Data.LocalSettings, "Local Settings (override)", w))
		b.WriteString("\n")
	}
	if ch.Data == nil || (ch.Data.Settings == nil && ch.Data.LocalSettings == nil) {
		b.WriteString(card("Settings", []string{
			emptyMsgStyle.Render("No settings.json found"),
			emptyHintStyle.Render("Configure in .claude/settings.json"),
		}, w))
		b.WriteString("\n")
	}

	if !ch.IsGlobal && ch.Data != nil {
		globalCh := m.findGlobalChannel()
		if globalCh != nil && globalCh.Data != nil && globalCh.Data.Settings != nil && ch.Data.Settings != nil {
			diffs := diffSettings(globalCh.Data.Settings, ch.Data.Settings)
			if len(diffs) > 0 {
				var lines []string
				for _, d := range diffs {
					val := fmt.Sprintf("%s → %s", labelStyle.Render(d.globalVal), valueStyle.Render(d.projectVal))
					lines = append(lines, cardKV(d.key, val, 12))
				}
				b.WriteString(card("Overrides from Global", lines, w))
				b.WriteString("\n")
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
		diffs = append(diffs, settingsDiff{"teammate", global.TeammateMode, project.TeammateMode})
	}
	return diffs
}

func renderSettingsCard(s *claude.Settings, title string, width int) string {
	var lines []string

	if s.Model != "" {
		lines = append(lines, cardKV("model", s.Model, 10))
	}
	if s.Language != "" {
		lines = append(lines, cardKV("language", s.Language, 10))
	}
	if s.TeammateMode != "" {
		lines = append(lines, cardKV("teammate", s.TeammateMode, 10))
	}
	if s.PlansDirectory != "" {
		lines = append(lines, cardKV("plans dir", s.PlansDirectory, 10))
	}

	result := card(title, lines, width)

	if len(s.Env) > 0 {
		var envLines []string
		for k, v := range s.Env {
			envLines = append(envLines, cardKV(k, v, 16))
		}
		result += "\n" + card("Environment", envLines, width)
	}

	if len(s.Permissions.Allow) > 0 || len(s.Permissions.Deny) > 0 {
		var permLines []string
		if len(s.Permissions.Allow) > 0 {
			permLines = append(permLines, sectionTitleStyle.Render("Allow"))
			for _, p := range s.Permissions.Allow {
				permLines = append(permLines, fmt.Sprintf("  %s %s", statusHealthy, p))
			}
		}
		if len(s.Permissions.Deny) > 0 {
			if len(s.Permissions.Allow) > 0 {
				permLines = append(permLines, "")
			}
			permLines = append(permLines, sectionTitleStyle.Render("Deny"))
			for _, p := range s.Permissions.Deny {
				permLines = append(permLines, fmt.Sprintf("  %s %s", statusError, p))
			}
		}
		result += "\n" + card("Permissions", permLines, width)
	}

	return result
}
