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
		fmt.Printf("%-30s  100%%[%s]\n", strings.Repeat(PercentChar, 30), filename)
	} else {
		fmt.Printf("%-30s  %d%%[%s]\r", strings.Repeat("#", 30*percent/100), percent, filename)
	}
}
