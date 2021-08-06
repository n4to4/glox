package lox

import (
	"log"
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return expr.TokenLiteral()
}

func (p AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
	w := &strings.Builder{}

	w.WriteString("(" + name)
	for _, exp := range exprs {
		w.WriteString(" ")
		s, ok := exp.Accept(p).(string)
		if !ok {
			log.Fatalf("not a string: %v", p)
		}
		w.WriteString(s)
	}
	w.WriteString(")")

	return w.String()
}
