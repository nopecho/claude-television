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
