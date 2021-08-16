package main

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	values := make(map[string]interface{}, 10)
	return &Environment{values}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (interface{}, error) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value, nil
	}

	return nil, RuntimeError{
		name,
		fmt.Sprintf("undefined variable %q", name.Lexeme),
	}
}

func (e *Environment) assign(name Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	return RuntimeError{name, fmt.Sprintf("undefined variable '%s'.", name.Lexeme)}
}
