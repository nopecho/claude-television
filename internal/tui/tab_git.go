package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderGitTab(ch *channel.Channel) string {
	var b strings.Builder

	if ch.Data.GitInfo == nil {
		b.WriteString(section("Git"))
		b.WriteString("    Not a git repository\n")
		return b.String()
	}

	git := ch.Data.GitInfo
	b.WriteString(section("Git"))
	b.WriteString(kv("branch", git.Branch) + "\n")

	if git.LastCommit != "" {
		b.WriteString(section("Last Commit"))
		b.WriteString(kv("hash", git.LastCommit) + "\n")
		b.WriteString(kv("message", git.LastCommitMsg) + "\n")
		b.WriteString(kv("date", git.LastCommitAt) + "\n")
	}

	b.WriteString(section("Working Tree"))
	if git.DirtyFiles > 0 {
		b.WriteString(fmt.Sprintf("    %s %d dirty files\n", statusWarning, git.DirtyFiles))
	} else {
		b.WriteString(fmt.Sprintf("    %s clean\n", statusHealthy))
	}
	return b.String()
}
