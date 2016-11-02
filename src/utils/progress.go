package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	PercentChar string = "#"
)

func PrintFileProgress(filename string, percent int) {

	if percent == 100 {
		time.Sleep(10 * time.Microsecond)
		fmt.Printf("\r%-50s  100%%[%s]\n", strings.Repeat(PercentChar, 50), filename)
	} else {
		fmt.Printf("\r%-50s  %d%%[%s]", strings.Repeat("#", percent*50/100), percent, filename)
	}
}
