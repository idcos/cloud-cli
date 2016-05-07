package util

import (
	"regexp"
	"strings"
)

// WildCharToRegexp change string with wildchar to regexp format string
func WildCharToRegexp(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "?", ".?", -1)
	s = strings.Replace(s, "*", ".*", -1)
	return "^" + s + "$"
}

func IsWildCharMatch(s string, wildCharStrs ...string) bool {

	var patterns = make([]string, 0)
	for _, s := range wildCharStrs {
		patterns = append(patterns, WildCharToRegexp(s))
	}

	for _, p := range patterns {
		matched, _ := regexp.MatchString(p, s)
		if matched {
			return true
		}
	}
	return false
}
