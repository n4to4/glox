package lox

import (
	"fmt"

	"github.com/n4to4/glox/tokens"
)

type Lox struct {
	//hadError bool
}

func LoxMain() {
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
	fmt.Println(p.Print(&expression))
}
