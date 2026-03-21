package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderHealthTab(ch *channel.Channel) string {
	var b strings.Builder

	issues := ch.Data.HealthIssues
	if len(issues) == 0 {
		b.WriteString(section("Health Check"))
		b.WriteString(fmt.Sprintf("    %s All checks passed\n", statusHealthy))
		return b.String()
	}

	errors := 0
	warnings := 0
	for _, i := range issues {
		switch i.Severity {
		case claude.SeverityError:
			errors++
		case claude.SeverityWarning:
			warnings++
		}
	}

	b.WriteString(section("Health Check"))
	if errors > 0 {
		b.WriteString(fmt.Sprintf("    %s %d errors  ", statusError, errors))
	}
	if warnings > 0 {
		b.WriteString(fmt.Sprintf("    %s %d warnings  ", statusWarning, warnings))
	}
	b.WriteString("\n")

	if errors > 0 {
		b.WriteString(section("Errors"))
		for _, i := range issues {
			if i.Severity == claude.SeverityError {
				b.WriteString(fmt.Sprintf("    %s %s\n", statusError, i.Message))
			}
		}
	}

	if warnings > 0 {
		b.WriteString(section("Warnings"))
		for _, i := range issues {
			if i.Severity == claude.SeverityWarning {
				b.WriteString(fmt.Sprintf("    %s %s\n", statusWarning, i.Message))
			}
		}
	}

	return b.String()
}
