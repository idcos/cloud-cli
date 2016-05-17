package util

import (
	"os/user"
	"regexp"
	"strings"
)

// HomeDirFlag 当前用户家目录标识符
const HomeDirFlag = "~"

// WildCharToRegexp change string with wildchar to regexp format string
func WildCharToRegexp(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, "?", ".?", -1)
	s = strings.Replace(s, "*", ".*", -1)
	return "^" + s + "$"
}

// IsWildCharMatch check s is match one of wildCharStrs or not
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

// ConvertHomeDir convert ~ to user's home dir
func ConvertHomeDir(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, HomeDirFlag) {
		return raw, nil
	}

	if !strings.HasPrefix(raw, HomeDirFlag) {
		return raw, nil
	}
	user, err := user.Current()
	if err != nil {
		return raw, err
	}
	return strings.Replace(raw, HomeDirFlag, user.HomeDir, 1), nil
}
