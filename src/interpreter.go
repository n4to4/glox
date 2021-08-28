package main

import (
	"fmt"
	"log"
	"time"
)

//
// native functions
//
type LoxClock struct{}

func (c LoxClock) Arity() int {
	return 0
}
func (c LoxClock) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	t := time.Now()
	ut := float64(t.UnixMilli()) / 1000.0
	return ut
}
func (c LoxClock) String() string {
	return "<native fn>"
}

//--------------------------------------------------------------------------------
// interpreter
//

type Interpreter struct {
	globals     *Environment
	environment *Environment
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.define("clock", LoxClock{})

	return &Interpreter{globals, globals}
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
	if stmt.initializer == nil {
		i.environment.define(stmt.name.lexeme, nil)
		return nil, nil
	}

	value, err := i.evaluate(*stmt.initializer)
	if err != nil {
		return nil, err
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt Block) (interface{}, error) {
	i.executeBlock(stmt.statements, NewEnvironment(i.environment))
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt If) (interface{}, error) {
	evaled, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(evaled) {
		return nil, i.execute(*stmt.thenBranch)
	}
	if stmt.elseBranch != nil {
		return nil, i.execute(*stmt.elseBranch)
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt While) (interface{}, error) {
	for {
		v, err := i.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(v) {
			break
		}

		if err := i.execute(stmt.body); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt Function) (interface{}, error) {
	function := LoxFunction{stmt, i.environment}
	i.environment.define(stmt.name.lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt Return) (interface{}, error) {
	var value interface{} = nil
	if stmt.value != nil {
		v, _ := i.evaluate(*stmt.value)
		value = v
	}

	return nil, ReturnValue{value}
}

func (i *Interpreter) execute(stmt Stmt) error {
	_, err := stmt.Accept(i)
	return err
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}

	return nil
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

	switch expr.operator.ttype {
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

	switch expr.operator.ttype {
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

func (i *Interpreter) VisitAssignExpr(expr Assign) (interface{}, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	if i.environment.assign(expr.name, value) != err {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr Logical) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	if expr.operator.ttype == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) VisitCallExpr(expr Call) (interface{}, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, argument := range expr.arguments {
		evaled, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, evaled)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, RuntimeError{
			expr.paren, "can only call functions and classes",
		}
	}
	if len(arguments) != function.Arity() {
		return nil, RuntimeError{
			expr.paren,
			fmt.Sprintf("expected %d arguments but got %d",
				function.Arity(), len(arguments)),
		}
	}
	return function.Call(i, arguments), nil
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
