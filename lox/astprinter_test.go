package lox

import (
	"testing"

	"github.com/n4to4/glox/tokens"
)

func TestAstPrinter(t *testing.T) {
	expression := Binary{
		Unary{
			tokens.NewToken(tokens.MINUS, "-", nil, 1),
			Literal{"123"},
		},
		tokens.NewToken(tokens.STAR, "*", nil, 1),
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
