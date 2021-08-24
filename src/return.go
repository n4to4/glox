package main

import "fmt"

type ReturnValue struct {
	value interface{}
}

func (r ReturnValue) Error() string {
	return fmt.Sprintf("<return %v>", r.value)
}
