package main

import (
	"fmt"
	"log"
)

type Interpreter struct{}

func (i *Interpreter) Interpret(expression Expr) {
	value, err := i.evaluate(expression)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(stringify(value))
}

func stringify(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	return expr.value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Ttype {
	case MINUS:
		return -(right.(float64)), nil
	case BANG:
		return !(isTruthy(right)), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Ttype {
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil

	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil

	case MINUS:
		return left.(float64) - right.(float64), nil
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r, nil
			}
		}
	case SLASH:
		return left.(float64) / right.(float64), nil
	case STAR:
		return left.(float64) * right.(float64), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func isTruthy(object interface{}) bool {
	switch v := object.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func isEqual(a, b interface{}) bool {
	return a == b
}
