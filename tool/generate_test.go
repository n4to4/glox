package main

/*
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

import "github.com/n4to4/glox/tokens"

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
	// pending
	if want != got {
		//t.Errorf("want %s, got %s", want, got)
	}
}
*/
