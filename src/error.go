package main

import (
	"fmt"
	"os"
)

func ErrorReport(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s", line, where, message)
}

func ReportError(token Token, message string) {
	if token.ttype == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}
