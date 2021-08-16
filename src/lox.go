package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	ErrOperandMustBeANumber     = "operand must be a number"
	ErrOperandsMustBeNumbers    = "operands must be numbers"
	ErrOperandsMustBeNumsOrStrs = "operands must be two numbers or two strings"
)

type RuntimeError struct {
	token   Token
	message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.line)
}

type Lox struct {
	interpreter     *Interpreter
	hadError        bool
	hadRuntimeError bool
}

func NewLox() Lox {
	return Lox{
		interpreter:     NewInterpreter(),
		hadError:        false,
		hadRuntimeError: false,
	}
}

func (l *Lox) runFile(file string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	l.run(string(bytes))
}

func (l *Lox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	prompt := func() { fmt.Print("> ") }
	for prompt(); scanner.Scan(); prompt() {
		line := scanner.Text()
		l.run(line)
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := Parser{tokens: tokens}
	stmts := parser.Parse()

	l.interpreter.Interpret(stmts)
}
