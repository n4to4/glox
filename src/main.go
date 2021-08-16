package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	}

	lox := NewLox()
	if len(os.Args) == 2 {
		lox.runFile(os.Args[1])
	} else {
		lox.runPrompt()
	}
}
