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
		return emptyState("Health Check", "All checks passed", "No issues found in project configuration")
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

	b.WriteString(section("Summary"))
	if errors > 0 {
		b.WriteString(sectionLine(fmt.Sprintf("  %s %d errors", statusError, errors)) + "\n")
	}
	if warnings > 0 {
		b.WriteString(sectionLine(fmt.Sprintf("  %s %d warnings", statusWarning, warnings)) + "\n")
	}

	if errors > 0 {
		b.WriteString(section("Errors"))
		for _, i := range issues {
			if i.Severity == claude.SeverityError {
				b.WriteString(sectionLine(fmt.Sprintf("  %s %s", statusError, i.Message)) + "\n")
			}
		}
	}

	if warnings > 0 {
		b.WriteString(section("Warnings"))
		for _, i := range issues {
			if i.Severity == claude.SeverityWarning {
				b.WriteString(sectionLine(fmt.Sprintf("  %s %s", statusWarning, i.Message)) + "\n")
			}
		}
	}

	return b.String()
}
