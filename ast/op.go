// Copyright (c) 2024 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"github.com/tsavola/dp/token"
)

// Special precedence levels.
const (
	MinBinaryPrecedence = 1
	MaxBinaryPrecedence = 5
	UltimatePrecedence  = 6
)

type BinaryOp token.Kind

const (
	OpAdd       = BinaryOp(token.Plus)
	OpSubtract  = BinaryOp(token.Minus)
	OpMultiply  = BinaryOp(token.Asterisk)
	OpDivide    = BinaryOp(token.Slash)
	OpRemainder = BinaryOp(token.Percent)

	OpLogicalAnd = BinaryOp(token.LogicalAnd)
	OpLogicalOr  = BinaryOp(token.LogicalOr)

	OpAndNot      = BinaryOp(token.AndNot)
	OpAnd         = BinaryOp(token.Ampersand)
	OpOr          = BinaryOp(token.Pipe)
	OpExclusiveOr = BinaryOp(token.Caret)

	OpShiftLeft  = BinaryOp(token.ShiftLeft)
	OpShiftRight = BinaryOp(token.ShiftRight)

	OpEqual          = BinaryOp(token.Equal)
	OpNotEqual       = BinaryOp(token.NotEqual)
	OpLessOrEqual    = BinaryOp(token.LessOrEqual)
	OpGreaterOrEqual = BinaryOp(token.GreaterOrEqual)
	OpLess           = BinaryOp(token.Less)
	OpGreater        = BinaryOp(token.Greater)
)

// Precedence level of the binary operator.
func (op BinaryOp) Precedence() int {
	switch op {
	case OpMultiply, OpDivide, OpRemainder, OpShiftLeft, OpShiftRight, OpAnd, OpAndNot:
		return 5

	case OpAdd, OpSubtract, OpOr, OpExclusiveOr:
		return 4

	case OpEqual, OpNotEqual, OpLess, OpLessOrEqual, OpGreater, OpGreaterOrEqual:
		return 3

	case OpLogicalAnd:
		return 2

	case OpLogicalOr:
		return 1
	}

	panic("invalid binary operator")
}

func (op BinaryOp) Dump() string {
	return "BinaryOp{" + op.String() + "}"
}

func (op BinaryOp) String() string {
	return token.Kind(op).String()
}

type UnaryOp token.Kind

const (
	OpIdentity   = UnaryOp(token.Plus)
	OpNegate     = UnaryOp(token.Minus)
	OpComplement = UnaryOp(token.Caret)
	OpNot        = UnaryOp(token.Exclamation)
)

func (op UnaryOp) Dump() string {
	return "UnaryOp{" + op.String() + "}"
}

func (op UnaryOp) String() string {
	return token.Kind(op).String()
}
