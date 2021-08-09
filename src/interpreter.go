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
		if err := checkNumberOperand(expr.operator, right); err != nil {
			return nil, err
		}
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
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil

	case MINUS:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
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
		return nil, RuntimeError{expr.operator, ErrOperandsMustBeNumsOrStrs}
	case SLASH:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case STAR:
		if err := checkNumberOperands(expr.operator, left, right); err != nil {
			return nil, err
		}
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

func stringify(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

func checkNumberOperand(operator Token, operand interface{}) error {
	_, ok := operand.(float64)
	if ok {
		return nil
	}

	return RuntimeError{operator, ErrOperandMustBeANumber}
}

func checkNumberOperands(operator Token, left, right interface{}) error {
	_, ok1 := left.(float64)
	_, ok2 := right.(float64)
	if ok1 && ok2 {
		return nil
	}

	return RuntimeError{operator, ErrOperandsMustBeNumbers}
}
