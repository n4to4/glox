package main

import (
	"fmt"
	"log"
)

type Interpreter struct {
	environment Environment
}

func NewInterpreter() Interpreter {
	env := NewEnvironment()
	return Interpreter{env}
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		if err := i.execute(stmt); err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}

//
// Visit Stmt
//

func (i *Interpreter) VisitExpressionStmt(stmt Expression) (interface{}, error) {
	return i.evaluate(stmt.expression)
}

func (i *Interpreter) VisitPrintStmt(stmt Print) (interface{}, error) {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(value)
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt Var) (interface{}, error) {
	fmt.Printf("name=%q, initializer=%v\n", stmt.name, stmt.initializer)

	if stmt.initializer == nil {
		i.environment.define(stmt.name.Lexeme, nil)
		return nil, nil
	}

	value, err := i.evaluate(*stmt.initializer)
	if err != nil {
		return nil, err
	}

	i.environment.define(stmt.name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) execute(stmt Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

//
// Visit Expr
//

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

func (i *Interpreter) VisitVariableExpr(expr Variable) (interface{}, error) {
	return i.environment.get(expr.name)
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
