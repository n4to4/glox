package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	values := make(map[string]interface{})
	return &Environment{enclosing, values}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (interface{}, error) {
	value, ok := e.values[name.lexeme]
	if ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, RuntimeError{
		name,
		fmt.Sprintf("undefined variable %q", name.lexeme),
	}
}

func (e *Environment) assign(name Token, value interface{}) error {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return RuntimeError{name, fmt.Sprintf("undefined variable '%s'.", name.lexeme)}
}
