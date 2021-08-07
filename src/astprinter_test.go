package main

import (
	"testing"
)

func TestAstPrinter(t *testing.T) {
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
	got := p.Print(expression)
	want := "(* (- 123) (group 45.67))"

	if want != got {
		t.Errorf("want %v got %v", want, got)
	}
}
