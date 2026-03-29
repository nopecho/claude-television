package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderHealthTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	issues := ch.Data.HealthIssues
	if len(issues) == 0 {
		b.WriteString(card("Health Check", []string{
			statusHealthy + " " + valueStyle.Render("All checks passed"),
			emptyHintStyle.Render("No issues found in project configuration"),
		}, w))
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

	// Summary card
	var summaryLines []string
	if errors > 0 {
		summaryLines = append(summaryLines, fmt.Sprintf("%s %d errors", statusError, errors))
	}
	if warnings > 0 {
		summaryLines = append(summaryLines, fmt.Sprintf("%s %d warnings", statusWarning, warnings))
	}
	b.WriteString(card("Summary", summaryLines, w))
	b.WriteString("\n")

	// Errors card
	if errors > 0 {
		var errLines []string
		for _, i := range issues {
			if i.Severity == claude.SeverityError {
				errLines = append(errLines, fmt.Sprintf("%s %s", statusError, i.Message))
			}
		}
		b.WriteString(card("Errors", errLines, w))
		b.WriteString("\n")
	}

	// Warnings card
	if warnings > 0 {
		var warnLines []string
		for _, i := range issues {
			if i.Severity == claude.SeverityWarning {
				warnLines = append(warnLines, fmt.Sprintf("%s %s", statusWarning, i.Message))
			}
		}
		b.WriteString(card("Warnings", warnLines, w))
		b.WriteString("\n")
	}

	return b.String()
}
