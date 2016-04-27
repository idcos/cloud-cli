package util

import "testing"

func TestWildCharToRegexp(t *testing.T) {
	var cases = []struct {
		inputStr  string
		expectStr string
	}{
		{"gr.up", `gr\.up`},
		{"gr*up", "gr.*up"},
		{"gr?up", "gr.?up"},
		{"gr.?up", "gr\\..?up"},
		{"gr.*up", "gr\\..*up"},
		{"gr.?.*up", "gr\\..?\\..*up"},
		{"gr.?*up", "gr\\..?.*up"},
	}

	for _, c := range cases {
		result := WildCharToRegexp(c.inputStr)
		if c.expectStr != result {
			t.Fatalf("expect: %s, result: %s", c.expectStr, result)
		}
	}
}
