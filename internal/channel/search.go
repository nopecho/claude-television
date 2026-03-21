package channel

import "github.com/sahilm/fuzzy"

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
