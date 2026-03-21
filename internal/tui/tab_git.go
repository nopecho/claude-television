package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderGitTab(ch *channel.Channel) string {
	var b strings.Builder

	if ch.Data.GitInfo == nil {
		return emptyState("Git", "Not a git repository", "")
	}

	git := ch.Data.GitInfo
	b.WriteString(section("Branch"))
	b.WriteString(sectionLine("  "+valueStyle.Render(git.Branch)) + "\n")

	if git.LastCommit != "" {
		b.WriteString(section("Last Commit"))
		b.WriteString(kv("hash", git.LastCommit, 8) + "\n")
		b.WriteString(kv("message", git.LastCommitMsg, 8) + "\n")
		b.WriteString(kv("date", git.LastCommitAt, 8) + "\n")
	}

	b.WriteString(section("Working Tree"))
	if git.DirtyFiles > 0 {
		b.WriteString(sectionLine(fmt.Sprintf("  %s %d dirty files", statusWarning, git.DirtyFiles)) + "\n")
	} else {
		b.WriteString(sectionLine(fmt.Sprintf("  %s clean", statusHealthy)) + "\n")
	}
	return b.String()
}
