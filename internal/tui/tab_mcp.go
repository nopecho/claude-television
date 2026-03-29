package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderMCPTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if len(ch.Data.MCPServers) == 0 {
		b.WriteString(card("MCP Servers", []string{
			emptyMsgStyle.Render("No MCP servers configured"),
			emptyHintStyle.Render("Configure in .claude/settings.json"),
		}, w))
		return b.String()
	}

	for _, s := range ch.Data.MCPServers {
		var lines []string
		source := labelStyle.Render(fmt.Sprintf("[%s]", s.Source))
		lines = append(lines, cardKV("type", s.Type, 8)+"  "+source)
		if s.Command != "" {
			lines = append(lines, cardKV("command", s.Command, 8))
		}
		if len(s.Args) > 0 {
			lines = append(lines, cardKV("args", strings.Join(s.Args, " "), 8))
		}
		if s.URL != "" {
			lines = append(lines, cardKV("url", s.URL, 8))
		}
		if len(s.Env) > 0 {
			lines = append(lines, "")
			lines = append(lines, labelStyle.Render("Environment:"))
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
				lines = append(lines, fmt.Sprintf("  %s = %s", labelStyle.Render(k), valueStyle.Render(display)))
			}
		}
		b.WriteString(card(s.Name, lines, w))
		b.WriteString("\n")
	}
	return b.String()
}
