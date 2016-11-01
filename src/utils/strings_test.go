package utils

import "testing"

func TestWildCharToRegexp(t *testing.T) {
	var cases = []struct {
		inputStr  string
		expectStr string
	}{
		{"gr.up", `^gr\.up$`},
		{"gr*up", "^gr.*up$"},
		{"gr?up", "^gr.?up$"},
		{"gr.?up", "^gr\\..?up$"},
		{"gr.*up", "^gr\\..*up$"},
		{"gr.?.*up", "^gr\\..?\\..*up$"},
		{"gr.?*up", "^gr\\..?.*up$"},
	}

	for _, c := range cases {
		result := WildCharToRegexp(c.inputStr)
		if c.expectStr != result {
			t.Fatalf("expect: %s, result: %s", c.expectStr, result)
		}
	}
}

func TestTrim(t *testing.T) {
	var cases = []struct {
		inputStr  string
		expectStr string
		cutSets   []string
	}{
		{" test ", "test", []string{" "}},
		{"\t test ", "test", []string{" ", "\t"}},
		{"\t\t\n test \n\t\n", "test", []string{" ", "\t", "\n"}},
		{"\t test \n\t", "test", []string{" ", "\t", "\n"}},
		{"\t test \t\n", "test", []string{" ", "\t", "\n"}},
		{"\n\t test\t \n", "test", []string{" ", "\t", "\n"}},
		{"test", "test", []string{}},
	}

	for _, c := range cases {
		result := Trim(c.inputStr, c.cutSets...)
		if c.expectStr != result {
			t.Errorf("expect: %s, result: %s", c.expectStr, result)
		}
	}
}
