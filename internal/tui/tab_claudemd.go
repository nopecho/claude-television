package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderClaudeMDTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if ch.Data.ClaudeMD != nil {
		md := ch.Data.ClaudeMD

		// Info card
		b.WriteString(card("CLAUDE.md", []string{
			cardKV("path", md.Path, 6),
			cardKV("lines", fmt.Sprintf("%d", md.LineCount), 6),
		}, w))
		b.WriteString("\n")

		if len(md.Sections) > 0 {
			var sectionLines []string
			for _, s := range md.Sections {
				sectionLines = append(sectionLines, "  • "+s)
			}
			b.WriteString(card("Sections", sectionLines, w))
			b.WriteString("\n")
		}

		// Content card
		var contentLines []string
		lines := strings.Split(md.Content, "\n")
		for _, l := range lines {
			contentLines = append(contentLines, l)
		}
		b.WriteString(card("Content", contentLines, w))
		b.WriteString("\n")
	} else {
		b.WriteString(card("CLAUDE.md", []string{
			emptyMsgStyle.Render("Not found"),
			emptyHintStyle.Render("Create CLAUDE.md in project root"),
		}, w))
		b.WriteString("\n")
	}

	if len(ch.Data.SubClaudeMDs) > 0 {
		var subLines []string
		for _, md := range ch.Data.SubClaudeMDs {
			subLines = append(subLines, sectionTitleStyle.Render(fmt.Sprintf("%s (%d lines)", md.Path, md.LineCount)))
			for _, s := range md.Sections {
				subLines = append(subLines, "  • "+s)
			}
			subLines = append(subLines, "")
		}
		b.WriteString(card(fmt.Sprintf("Sub CLAUDE.md (%d)", len(ch.Data.SubClaudeMDs)), subLines, w))
		b.WriteString("\n")
	}
	return b.String()
}
