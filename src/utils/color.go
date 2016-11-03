package utils

import "github.com/fatih/color"

var (
	FgYellow = color.New(color.FgYellow).SprintFunc()
	FgBlue   = color.New(color.FgBlue).SprintFunc()
	FgRed    = color.New(color.FgRed).SprintFunc()
	FgCyan   = color.New(color.FgCyan).SprintFunc()
	FgGreen  = color.New(color.FgGreen).SprintFunc()

	FgBoldYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	FgBoldBlue   = color.New(color.FgBlue, color.Bold).SprintFunc()
	FgBoldRed    = color.New(color.FgRed, color.Bold).SprintFunc()
	FgBoldCyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	FgBoldGreen  = color.New(color.FgGreen, color.Bold).SprintFunc()

	BgYellow = color.New(color.BgYellow).SprintFunc()
	BgBlue   = color.New(color.BgBlue).SprintFunc()
	BgRed    = color.New(color.BgRed).SprintFunc()
	BgCyan   = color.New(color.BgCyan).SprintFunc()
	BgGreen  = color.New(color.BgGreen).SprintFunc()

	BgBoldYellow = color.New(color.BgYellow, color.Bold).SprintFunc()
	BgBoldBlue   = color.New(color.BgBlue, color.Bold).SprintFunc()
	BgBoldRed    = color.New(color.BgRed, color.Bold).SprintFunc()
	BgBoldCyan   = color.New(color.BgCyan, color.Bold).SprintFunc()
	BgBoldGreen  = color.New(color.BgGreen, color.Bold).SprintFunc()
)
