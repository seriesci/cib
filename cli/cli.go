package cli

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	check = "\u2713"
)

// cli colors
var (
	Green = color.New(color.FgGreen).SprintFunc()
	Blue  = color.New(color.FgBlue).SprintFunc()
)

// Checkln writes to standard output.
func Checkln(a ...interface{}) {
	fmt.Println(append([]interface{}{Green(check), "cib:"}, a...)...)
}

// Checkf formats according to a format specifier and writes to standard output.
func Checkf(format string, a ...interface{}) {
	// create new a with content all blue
	na := make([]interface{}, len(a))
	for index, content := range a {
		na[index] = Blue(content)
	}
	fmt.Printf(Green(check)+" cib: "+format, na...)
}
