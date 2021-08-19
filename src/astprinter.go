package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) Print(expr Expr) (string, error) {
	ret, err := expr.Accept(p)
	if err != nil {
		return "", err
	}

	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("not a string: %v", ret)
	}

	return str, nil
}

func (p AstPrinter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p AstPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (p AstPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p AstPrinter) VisitVariableExpr(expr Variable) (interface{}, error) {
	return nil, nil
}

func (p AstPrinter) VisitAssignExpr(expr Assign) (interface{}, error) {
	return nil, nil
}

func (p AstPrinter) VisitLogicalExpr(expr Logical) (interface{}, error) {
	return nil, nil
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	w := &strings.Builder{}

	w.WriteString("(" + name)
	for _, exp := range exprs {
		w.WriteString(" ")

		ret, err := exp.Accept(p)
		if err != nil {
			return "", err
		}

		s, ok := ret.(string)
		if !ok {
			return "", fmt.Errorf("not a string: %v", p)
		}

		w.WriteString(s)
	}
	w.WriteString(")")

	return w.String(), nil
}
