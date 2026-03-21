package util

import (
	"os"
	"strings"
)

// DecodeProjectPath decodes an encoded project directory name back to an absolute path.
// Claude Code encodes paths like "/Users/foo/projects/bar" as "-Users-foo-projects-bar".
func DecodeProjectPath(encoded string) string {
	if !strings.HasPrefix(encoded, "-") {
		return encoded
	}
	parts := strings.Split(strings.TrimPrefix(encoded, "-"), "-")
	return bestEffortDecode(parts)
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
