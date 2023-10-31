package service

import (
	"fmt"
)

type SyntaxError struct {
	msg string
}

func (e SyntaxError) Error() string {
	return e.msg
}

func NewSyntaxError(msg string) SyntaxError {
	return SyntaxError{msg: msg}
}

func NewSyntaxErrorf(format string, args ...any) SyntaxError {
	return NewSyntaxError(fmt.Sprintf(format, args...))
}
