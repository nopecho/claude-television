package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderGitTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if ch.Data.GitInfo == nil {
		b.WriteString(card("Git", []string{
			emptyMsgStyle.Render("Not a git repository"),
		}, w))
		return b.String()
	}

	git := ch.Data.GitInfo

	// Branch card
	b.WriteString(card("Branch", []string{
		valueStyle.Render(git.Branch),
	}, w))
	b.WriteString("\n")

	// Last Commit card
	if git.LastCommit != "" {
		b.WriteString(card("Last Commit", []string{
			cardKV("hash", git.LastCommit, 8),
			cardKV("message", git.LastCommitMsg, 8),
			cardKV("date", git.LastCommitAt, 8),
		}, w))
		b.WriteString("\n")
	}

	// Working Tree card
	var treeLines []string
	if git.DirtyFiles > 0 {
		treeLines = append(treeLines, fmt.Sprintf("%s %d dirty files", statusWarning, git.DirtyFiles))
	} else {
		treeLines = append(treeLines, fmt.Sprintf("%s clean", statusHealthy))
	}
	b.WriteString(card("Working Tree", treeLines, w))

	return b.String()
}
