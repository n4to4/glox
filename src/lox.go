package main

import (
	"fmt"
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
