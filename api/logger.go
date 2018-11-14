package api

import (
	"fmt"

	"github.com/fatih/color"
)

// Info prints info log
func Info(msg string) {
	fmt.Printf("%s %s\n", color.New(color.FgYellow, color.Bold).SprintFunc()("INFO"), msg)
}

// Error prints error log
func Error(msg string) {
	fmt.Printf("%s %s\n", color.New(color.FgRed, color.Bold).SprintFunc()("ERROR"), msg)
}
