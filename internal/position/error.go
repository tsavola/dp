// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package position

import (
	"fmt"

	"github.com/tsavola/dp/source"
)

type Error interface {
	error
	PositionError() string
}

type posError struct {
	pos  source.Position
	msg  string
	errs []error
}

func NewError(pos source.Position, msg string, errs ...error) error {
	return posError{pos, msg, errs}
}

func Errorf(pos source.Position, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	msg := err.Error()

	var wrapped []error

	switch wrapper := err.(type) {
	case interface{ Unwrap() []error }:
		wrapped = wrapper.Unwrap()

	case interface{ Unwrap() error }:
		if e := wrapper.Unwrap(); e != nil {
			wrapped = []error{e}
		}

	default:
	}

	return posError{pos, msg, wrapped}
}

func (e posError) Pos() source.Position  { return e.pos }
func (e posError) Error() string         { return e.msg }
func (e posError) PositionError() string { return e.IndentError("") }
func (e posError) Unwrap() []error       { return e.errs }

func (e posError) IndentError(indent string) string {
	var s string
	if e.pos.Path != "" {
		s += e.pos.Path + ":"
	}
	if e.pos.Line > 0 {
		s += fmt.Sprintf("%04d:%03d: ", e.pos.Line, e.pos.Column)
	}
	s += indent + e.msg

	if len(e.errs) > 0 {
		s += ":"
		indent = "  " + indent
		for _, e := range e.errs {
			s += "\n" + indentError(indent, e)
		}
	}

	return s
}

type errorIndenter interface {
	IndentError(indent string) string
}

func indentError(indent string, err error) string {
	if e, ok := err.(errorIndenter); ok {
		return e.IndentError(indent)
	} else {
		return indent + err.Error()
	}
}
