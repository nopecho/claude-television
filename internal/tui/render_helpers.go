package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func statusIconStr(s channel.ChannelStatus) string {
	switch s {
	case channel.StatusHealthy:
		return statusHealthy
	case channel.StatusWarning:
		return statusWarning
	case channel.StatusError:
		return statusError
	}
	return "?"
}

func boolIcon(b bool) string {
	if b {
		return statusHealthy
	}
	return statusError
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func kv(key, value string) string {
	return fmt.Sprintf("  %s  %s", labelStyle.Render(key+":"), valueStyle.Render(value))
}

func section(title string) string {
	return "\n" + headerStyle.Render("  "+title) + "\n"
}

func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "    " + l
	}
	return strings.Join(lines, "\n")
}

func bullet(s string) string {
	return "    • " + s
}
