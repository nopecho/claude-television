package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderClaudeMDTab(ch *channel.Channel) string {
	var b strings.Builder

	if ch.Data.ClaudeMD != nil {
		md := ch.Data.ClaudeMD
		b.WriteString(section("CLAUDE.md"))
		b.WriteString(kv("path", md.Path) + "\n")
		b.WriteString(kv("lines", fmt.Sprintf("%d", md.LineCount)) + "\n")
		if len(md.Sections) > 0 {
			b.WriteString(section("Sections"))
			for _, s := range md.Sections {
				b.WriteString(fmt.Sprintf("    • %s\n", s))
			}
		}
		b.WriteString(section("Content"))
		lines := strings.Split(md.Content, "\n")
		for _, l := range lines {
			b.WriteString("    " + l + "\n")
		}
	} else {
		b.WriteString(section("CLAUDE.md"))
		b.WriteString("    Not found\n")
	}

	if len(ch.Data.SubClaudeMDs) > 0 {
		b.WriteString(section(fmt.Sprintf("Sub CLAUDE.md files (%d)", len(ch.Data.SubClaudeMDs))))
		for _, md := range ch.Data.SubClaudeMDs {
			b.WriteString(fmt.Sprintf("    %s (%d lines)\n", md.Path, md.LineCount))
			for _, s := range md.Sections {
				b.WriteString(fmt.Sprintf("      • %s\n", s))
			}
		}
	}
	return b.String()
}
