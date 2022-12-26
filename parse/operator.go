// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
	"import.name/pan"
)

func parsePrefixOperator(s scan) (scan, ast.UnaryOp) {
	switch t := s.peek().Kind; t {
	case token.Plus, token.Minus, token.Caret, token.Exclamation:
		s.skim(t)
		return s, ast.UnaryOp(t)

	default:
		panic(pan.Wrap(newError(s.pos(), "prefix operator expected")))
	}
}

func parseInfixOperator(s scan) (scan, ast.BinaryOp, bool) {
	switch t := s.peek().Kind; t {
	case token.Plus, token.Minus, token.Asterisk, token.Slash, token.Percent,
		token.LogicalAnd, token.LogicalOr,
		token.AndNot, token.Ampersand, token.Pipe, token.Caret,
		token.ShiftLeft, token.ShiftRight,
		token.Equal, token.NotEqual, token.LessOrEqual, token.GreaterOrEqual, token.Less, token.Greater:
		s.skim(t)
		return s, ast.BinaryOp(t), true

	default:
		return s, 0, false
	}
}
