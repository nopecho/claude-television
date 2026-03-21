package tui

import (
	"fmt"
	"strings"
)

func (m model) buildSkillItems() []skillItem {
	var items []skillItem
	for i, p := range m.data.Plugins {
		icon := statusIconOff
		if p.Enabled {
			icon = statusIcon
		}
		tag := "plugin"
		if !p.Installed {
			tag = "not installed"
		}
		items = append(items, skillItem{
			display:  fmt.Sprintf("%s %s [%s]", icon, p.Name, tag),
			isPlugin: true, pluginIdx: i,
		})
	}
	for i, s := range m.data.LocalSkills {
		items = append(items, skillItem{
			display:  fmt.Sprintf("%s %s [local]", statusIcon, s.Name),
			isPlugin: false, skillIdx: i,
		})
	}
	return items
}

func (m model) renderSkillsList() string {
	if len(m.skillCache) == 0 {
		return listItemStyle.Render("No plugins or skills found.")
	}
	var b strings.Builder
	for i, item := range m.skillCache {
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + item.display))
		} else {
			b.WriteString(listItemStyle.Render("  " + item.display))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderSkillsDetail() string {
	if len(m.skillCache) == 0 || m.cursor >= len(m.skillCache) {
		return ""
	}
	item := m.skillCache[m.cursor]
	var b strings.Builder
	if item.isPlugin {
		p := m.data.Plugins[item.pluginIdx]
		b.WriteString(titleStyle.Render(p.Name) + "\n\n")
		b.WriteString(fmt.Sprintf("  key:         %s\n", p.Key))
		b.WriteString(fmt.Sprintf("  marketplace: %s\n", p.Marketplace))
		b.WriteString(fmt.Sprintf("  version:     %s\n", p.Version))
		b.WriteString(fmt.Sprintf("  enabled:     %s\n", boolIcon(p.Enabled)))
		b.WriteString(fmt.Sprintf("  installed:   %s\n", boolIcon(p.Installed)))
		if p.InstallPath != "" {
			b.WriteString(fmt.Sprintf("  path:        %s\n", p.InstallPath))
		}
	} else {
		s := m.data.LocalSkills[item.skillIdx]
		b.WriteString(titleStyle.Render(s.Name) + "\n\n")
		b.WriteString("  type: local skill\n")
		b.WriteString(fmt.Sprintf("  path: %s\n", s.Path))
	}
	return b.String()
}
