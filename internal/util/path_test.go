package util

import "testing"

func TestDecodeProjectPath_NoPrefix(t *testing.T) {
	got := DecodeProjectPath("plain-path")
	if got != "plain-path" {
		t.Errorf("expected plain-path, got %s", got)
	}
}

func TestDecodeProjectPath_RootPrefix(t *testing.T) {
	got := DecodeProjectPath("-tmp")
	if got != "/tmp" {
		t.Errorf("expected /tmp, got %s", got)
	}
}

func TestDecodeProjectPath_EmptyParts(t *testing.T) {
	got := DecodeProjectPath("-")
	if got != "/" {
		t.Errorf("expected /, got %s", got)
	}
}

func TestPreprocessDotParts(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "no dots",
			input: []string{"Users", "foo", "projects"},
			want:  []string{"Users", "foo", "projects"},
		},
		{
			name:  "single dot prefix",
			input: []string{"Users", "foo", "", "config"},
			want:  []string{"Users", "foo", ".config"},
		},
		{
			name:  "dot prefix mid path",
			input: []string{"Users", "foo", "", "local", "share", "chezmoi"},
			want:  []string{"Users", "foo", ".local", "share", "chezmoi"},
		},
		{
			name:  "multiple dot prefixes",
			input: []string{"Users", "foo", "", "config", "", "nvim"},
			want:  []string{"Users", "foo", ".config", ".nvim"},
		},
		{
			name:  "dot at home level",
			input: []string{"Users", "foo", "", "claude"},
			want:  []string{"Users", "foo", ".claude"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := preprocessDotParts(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("len mismatch: got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("part[%d]: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDecodeProjectPath_DotPrefixed(t *testing.T) {
	// "-Users-nopecho--config" should decode with ".config" in the path
	got := DecodeProjectPath("-Users-nopecho--config")
	// We can't assert the exact path since it depends on filesystem,
	// but we can verify it contains .config
	if got == "" {
		t.Error("expected non-empty path")
	}
	// The decoded path should contain ".config" not just "config"
	if !contains(got, ".config") {
		t.Errorf("expected .config in path, got %s", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
