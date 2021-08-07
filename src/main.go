package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(file string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	run(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		run(line)
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	fmt.Println("print tokens: ")
	for _, token := range tokens {
		fmt.Println(token)
	}

	parser := Parser{tokens: tokens}
	expr := parser.Parse()

	printer := AstPrinter{}

	fmt.Println("\nprint expr:")
	fmt.Println(printer.Print(expr))
}
