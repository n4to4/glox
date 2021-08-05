package main

import (
	"strings"
	"testing"
)

func TestGenerateAst(t *testing.T) {
	w := &strings.Builder{}
	generateAst(w, "Expr", []string{
		"Binary   : left Expr, operator Token, right Expr",
		//"Grouping : expression Expr",
		//"Literal  : value interface{}",
		//"Unary    : operator Token, right Expr",
	})

	got := w.String()
	want := `package lox

type Expr interface {
	TokenLiteral() string
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func (x *Binary) TokenLiteral() string { return "" }
`
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}
