package channel

import "testing"

func TestFuzzySearch(t *testing.T) {
	channels := []Channel{
		{ID: "1", Name: "my-api-server", Path: "/projects/my-api-server", Group: "work"},
		{ID: "2", Name: "web-frontend", Path: "/projects/web-frontend", Group: "work"},
		{ID: "3", Name: "personal-blog", Path: "/home/blog", Group: "personal"},
		{ID: "4", Name: "api-gateway", Path: "/projects/api-gateway"},
	}

	results := FuzzySearch(channels, "api")
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'api', got %d", len(results))
	}

	results = FuzzySearch(channels, "web")
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'web', got %d", len(results))
	}

	results = FuzzySearch(channels, "")
	if len(results) != 4 {
		t.Errorf("expected all 4 for empty query, got %d", len(results))
	}
}
