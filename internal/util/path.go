package util

import (
	"os"
	"strings"
)

// DecodeProjectPath decodes an encoded project directory name back to an absolute path.
// Claude Code encodes paths by replacing "/" with "-" and "." with "-".
// Examples:
//   - "-Users-foo-projects-bar" → "/Users/foo/projects/bar"
//   - "-Users-foo--config"     → "/Users/foo/.config"  (double dash = dot-prefixed)
//   - "-Users-foo--local-share-chezmoi" → "/Users/foo/.local/share/chezmoi"
func DecodeProjectPath(encoded string) string {
	if !strings.HasPrefix(encoded, "-") {
		return encoded
	}
	// Split on "-" and preprocess: an empty part means the next part was dot-prefixed.
	// e.g. "-Users-foo--config" → ["Users", "foo", "", "config"]
	//   → ["Users", "foo", ".config"]
	raw := strings.Split(strings.TrimPrefix(encoded, "-"), "-")
	parts := preprocessDotParts(raw)
	return bestEffortDecode(parts)
}

// preprocessDotParts merges empty parts with the following part as a dot prefix.
// ["Users", "foo", "", "config"] → ["Users", "foo", ".config"]
// ["Users", "foo", "", "local", "share"] → ["Users", "foo", ".local", "share"]
func preprocessDotParts(raw []string) []string {
	parts := make([]string, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		if raw[i] == "" && i+1 < len(raw) {
			parts = append(parts, "."+raw[i+1])
			i++ // skip next
		} else if raw[i] != "" {
			parts = append(parts, raw[i])
		}
	}
	return parts
}

func bestEffortDecode(parts []string) string {
	if len(parts) == 0 {
		return "/"
	}
	current := "/"
	i := 0
	for i < len(parts) {
		found := false
		for j := len(parts); j > i; j-- {
			candidate := current + strings.Join(parts[i:j], "-")
			if _, err := os.Stat(candidate); err == nil {
				current = candidate + "/"
				i = j
				found = true
				break
			}
		}
		if !found {
			current += parts[i] + "/"
			i++
		}
	}
	return strings.TrimSuffix(current, "/")
}
