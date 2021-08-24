package main

import (
	"errors"
	"fmt"
)

type LoxFunction struct {
	declaration Function
}

func (f LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnvironment(interpreter.globals)
	for i := 0; i < len(f.declaration.params); i++ {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.body, environment)
	var returnValue ReturnValue
	if errors.As(err, &returnValue) {
		return returnValue.value
	}

	return nil
}

func (f LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}
