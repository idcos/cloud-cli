package util

import (
	"strings"
)

// WildCharToRegexp change string with wildchar to regexp format string
func WildCharToRegexp(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "?", ".?", -1)
	s = strings.Replace(s, "*", ".*", -1)
	return s
}
