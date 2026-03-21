package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderMCPTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.MCPServers) == 0 {
		return emptyState("MCP Servers", "No MCP servers configured", "Configure in .claude/settings.json")
	}

	b.WriteString(section(fmt.Sprintf("MCP Servers (%d)", len(ch.Data.MCPServers))))

	for _, s := range ch.Data.MCPServers {
		source := labelStyle.Render(fmt.Sprintf("[%s]", s.Source))
		b.WriteString(sectionEmpty() + "\n")
		b.WriteString(sectionLine(sectionTitleStyle.Render(s.Name)+" "+source) + "\n")
		b.WriteString(sectionLine(fmt.Sprintf("  type: %s", s.Type)) + "\n")
		if s.Command != "" {
			b.WriteString(sectionLine(fmt.Sprintf("  command: %s", s.Command)) + "\n")
		}
		if len(s.Args) > 0 {
			b.WriteString(sectionLine(fmt.Sprintf("  args: %s", strings.Join(s.Args, " "))) + "\n")
		}
		if s.URL != "" {
			b.WriteString(sectionLine(fmt.Sprintf("  url: %s", s.URL)) + "\n")
		}
		if len(s.Env) > 0 {
			b.WriteString(sectionLine("  env:") + "\n")
			for k, v := range s.Env {
				display := v
				lower := strings.ToLower(k)
				if strings.Contains(lower, "key") || strings.Contains(lower, "secret") || strings.Contains(lower, "token") {
					if len(v) > 4 {
						display = v[:4] + "****"
					} else {
						display = "****"
					}
				}
				b.WriteString(sectionLine(fmt.Sprintf("    %s = %s", k, display)) + "\n")
			}
		}
	}
	return b.String()
}
