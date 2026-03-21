package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderMCPTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.MCPServers) == 0 {
		b.WriteString(section("MCP Servers"))
		b.WriteString("    No MCP servers configured\n")
		return b.String()
	}

	b.WriteString(section(fmt.Sprintf("MCP Servers (%d)", len(ch.Data.MCPServers))))

	for _, s := range ch.Data.MCPServers {
		source := labelStyle.Render(fmt.Sprintf("[%s]", s.Source))
		b.WriteString(fmt.Sprintf("\n    %s %s\n", headerStyle.Render(s.Name), source))
		b.WriteString(kv("    type", s.Type) + "\n")
		if s.Command != "" {
			b.WriteString(kv("    command", s.Command) + "\n")
		}
		if len(s.Args) > 0 {
			b.WriteString(kv("    args", strings.Join(s.Args, " ")) + "\n")
		}
		if s.URL != "" {
			b.WriteString(kv("    url", s.URL) + "\n")
		}
		if len(s.Env) > 0 {
			b.WriteString("      env:\n")
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
				b.WriteString(fmt.Sprintf("        %s = %s\n", k, display))
			}
		}
	}
	return b.String()
}
