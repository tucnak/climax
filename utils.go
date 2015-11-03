package climax

import (
	"fmt"
	"io"
	"os"
)

var (
	outputDevice io.Writer = os.Stdout
	errorDevice  io.Writer = os.Stderr
)

func println(stuff ...interface{}) {
	fmt.Fprintln(outputDevice, stuff...)
}

func printf(format string, stuff ...interface{}) {
	fmt.Fprintf(outputDevice, format, stuff...)
}

func printerr(err error) {
	fmt.Fprintln(errorDevice, err)
}
