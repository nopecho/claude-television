package channel

import (
	"strings"

	"github.com/sahilm/fuzzy"
)

type searchable struct {
	channels []Channel
}

func (s searchable) String(i int) string {
	ch := s.channels[i]
	return ch.Name + " " + ch.Path + " " + ch.Group
}

func (s searchable) Len() int {
	return len(s.channels)
}

// ContentSearch searches across channel data content (settings, MCP, hooks, plugins).
func ContentSearch(channels []Channel, query string) []Channel {
	if query == "" {
		return channels
	}
	q := strings.ToLower(query)
	var results []Channel
	for _, ch := range channels {
		if ch.Data == nil {
			continue
		}
		if matchesContent(ch.Data, q) {
			results = append(results, ch)
		}
	}
	return results
}

func matchesContent(data *ChannelData, query string) bool {
	if data.Settings != nil {
		if containsLower(data.Settings.Model, query) ||
			containsLower(data.Settings.Language, query) ||
			containsLower(data.Settings.TeammateMode, query) {
			return true
		}
		for _, p := range data.Settings.Permissions.Allow {
			if containsLower(p, query) {
				return true
			}
		}
	}
	for _, s := range data.MCPServers {
		if containsLower(s.Name, query) || containsLower(s.Command, query) {
			return true
		}
	}
	for _, h := range data.Hooks {
		if containsLower(h.Event, query) || containsLower(h.Command, query) || containsLower(h.Matcher, query) {
			return true
		}
	}
	for _, p := range data.Plugins {
		if containsLower(p.Name, query) {
			return true
		}
	}
	for _, s := range data.LocalSkills {
		if containsLower(s.Name, query) {
			return true
		}
	}
	return false
}

func containsLower(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}

func FuzzySearch(channels []Channel, query string) []Channel {
	if query == "" {
		return channels
	}
	source := searchable{channels: channels}
	matches := fuzzy.FindFrom(query, source)
	result := make([]Channel, 0, len(matches))
	for _, m := range matches {
		result = append(result, channels[m.Index])
	}
	return result
}
