package utils

import (
	"fmt"
	"strings"
)

const (
	PercentChar string = "#"
)

func PrintFileProgress(filename string, percent int) {

	if percent == 100 {
		fmt.Printf("%-100s  100%%[%s]\n", strings.Repeat(PercentChar, 100), filename)
	} else {
		fmt.Printf("%-100s  %d%%[%s]\r", strings.Repeat("#", percent), percent, filename)
	}
}
