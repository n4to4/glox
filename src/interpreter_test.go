package main

import (
	"errors"
	"reflect"
	"testing"
)

var (
	minus = NewToken(MINUS, "-", nil, 1)
	plus  = NewToken(PLUS, "+", nil, 2)
)

func TestInterpreter(t *testing.T) {
	cases := []struct {
		expr Expr
		want interface{}
	}{
		{
			expr: Literal{1},
			want: 1,
		},
		{
			expr: Binary{Literal{"abc"}, plus, Literal{"123"}},
			want: "abc123",
		},
	}

	interpreter := Interpreter{}
	for _, cc := range cases {
		got, err := interpreter.evaluate(cc.expr)

		if err != nil {
			t.Errorf("want no error, got %v", err)
		}

		if !reflect.DeepEqual(got, cc.want) {
			t.Errorf("want %v, got %v", cc.want, got)
		}
	}
}

func TestInterpreterError(t *testing.T) {
	cases := []struct {
		expr Expr
		want error
	}{
		{
			expr: Unary{minus, Literal{"string"}},
			want: RuntimeError{minus, ErrOperandMustBeANumber},
		},
	}

	interpreter := Interpreter{}
	for _, cc := range cases {
		_, err := interpreter.evaluate(cc.expr)

		if err == nil {
			t.Error("want error, got none")
		}

		if !errors.Is(err, cc.want) {
			t.Errorf("want error %v, got %v", cc.want, err)
		}
	}
}
