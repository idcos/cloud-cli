package util

import (
	"fmt"
	"strings"
)

// Confirm 标准输入确认
func Confirm(prompt string) bool {
	for {
		var input string
		fmt.Printf(prompt)
		fmt.Scanln(&input)
		input = strings.ToLower(input)

		switch input {
		case "no", "n":
			return false
		case "yes", "y":
			return true
		}
	}
}
