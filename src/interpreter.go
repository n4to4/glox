package main

import "fmt"

type Interpreter struct{}

func (i *Interpreter) Interpret(expression Expr) {
	value := i.evaluate(expression)
	fmt.Println(stringify(value))
}

func stringify(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.Ttype {
	case MINUS:
		return -(right.(float64))
	case BANG:
		return !(isTruthy(right))
	}

	// unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.Ttype {
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)

	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)

	case MINUS:
		return left.(float64) - right.(float64)
	case PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	}

	// unreachable
	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
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
