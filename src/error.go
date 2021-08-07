package main

import (
	"fmt"
	"os"
)

// error

func ErrorReport(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s", line, where, message)
}
