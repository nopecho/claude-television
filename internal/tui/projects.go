package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderProjectsList() string {
	if len(m.data.Projects) == 0 {
		return listItemStyle.Render("No projects found.\nUse: ctv scan <path>")
	}
	var b strings.Builder
	for i, p := range m.data.Projects {
		icon := statusIconOff
		if p.HasClaudeDir || p.HasClaudeMD {
			icon = statusIcon
		}
		line := fmt.Sprintf("%s %s", icon, p.Name)
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + line))
		} else {
			b.WriteString(listItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderProjectsDetail() string {
	if len(m.data.Projects) == 0 || m.cursor >= len(m.data.Projects) {
		return ""
	}
	p := m.data.Projects[m.cursor]
	var b strings.Builder
	b.WriteString(titleStyle.Render(p.Name) + "\n")
	b.WriteString(fmt.Sprintf("  %s\n\n", p.Path))
	b.WriteString(fmt.Sprintf("  .claude/          %s\n", boolIcon(p.HasClaudeDir)))
	b.WriteString(fmt.Sprintf("  CLAUDE.md         %s\n", boolIcon(p.HasClaudeMD)))
	b.WriteString(fmt.Sprintf("  settings.json     %s\n", boolIcon(p.HasSettings)))

	meta := m.findProjectMeta(p.Path)
	if meta != nil {
		b.WriteString(fmt.Sprintf("\n  memory            %s\n", boolIcon(meta.HasMemory)))
		b.WriteString(fmt.Sprintf("  sessions          %s\n", boolIcon(meta.HasSessions)))
	}

	if p.HasClaudeMD {
		md, err := claude.ParseClaudeMD(p.Path + "/CLAUDE.md")
		if err == nil {
			b.WriteString(fmt.Sprintf("\n  CLAUDE.md (%d lines):\n", md.LineCount))
			for _, s := range md.Sections {
				b.WriteString(fmt.Sprintf("    • %s\n", s))
			}
		}
	}
	return b.String()
}

func (m model) findProjectMeta(path string) *claude.ProjectMeta {
	for i := range m.data.ProjectsMeta {
		if m.data.ProjectsMeta[i].DecodedPath == path {
			return &m.data.ProjectsMeta[i]
		}
	}
	return nil
}

func boolIcon(b bool) string {
	if b {
		return statusIcon
	}
	return statusIconOff
}
