package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Confirm 标准输入确认
func Confirm(promptStr string) bool {
	for {
		input := strings.ToLower(prompt(promptStr))

		switch input {
		case "no", "n":
			return false
		case "yes", "y":
			return true
		}
	}
}

// LoginNo which node you want to login [min, max)
func LoginNo(promptStr string, min, max int) int {
	for {
		input := prompt(promptStr)
		i, _ := strconv.Atoi(input)
		if i < min || i >= max {
			continue
		}

		return i
	}
}

func prompt(promptStr string) string {
	var input string
	fmt.Printf(promptStr)
	fmt.Scanln(&input)
	return input
}
