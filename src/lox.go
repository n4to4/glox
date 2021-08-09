package main

import (
	"fmt"
)

const (
	ErrOperandMustBeANumber     = "operand must be a number"
	ErrOperandsMustBeNumbers    = "operands must be numbers"
	ErrOperandsMustBeNumsOrStrs = "operands must be two numbers or two strings"
)

type Lox struct {
	//hadError bool
}

func LoxMain() {
	expression := Binary{
		Unary{
			NewToken(MINUS, "-", nil, 1),
			Literal{"123"},
		},
		NewToken(STAR, "*", nil, 1),
		Grouping{
			Literal{"45.67"},
		},
	}

	p := AstPrinter{}
	fmt.Println(p.Print(&expression))
}

type RuntimeError struct {
	token   Token
	message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", e.message, e.token.line)
}
