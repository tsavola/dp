// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
	"import.name/pan"
)

func parseExpressionInBlock(s scan) (scan, ast.BlockChild) {
	s, expr := parseAnyExpr(s, false)

	switch s.peek().Kind {
	case token.Comment, token.Newline, token.Semicolon:
	default:
		pan.Panic(newError(s.pos(), "expression: end of statement expected"))
	}

	return s, ast.Expression{expr}
}

func parseExpressionInExprList(s scan) (scan, ast.ExprListChild) {
	s, node, ok := parseAssignerDereference(s)
	if ok {
		return s, node
	}

	s, expr := parseAnyExpr(s, false)

	if !s.skip(token.Comma) {
		switch s.peek().Kind {
		case token.BraceRight, token.Comment, token.Newline, token.ParenRight, token.Semicolon:
		default:
			pan.Panic(newError(s.pos(), "end of expression expected"))
		}
	}

	return s, ast.Expression{expr}
}

func parseAnyExpr(s scan, multiline bool) (scan, ast.ExprChild) {
	var left ast.ExprChild
	var operator ast.BinaryOp
	var secondary bool

	for {
		if multiline {
			for s.skip(token.Newline) {
			}
		}

		var operand ast.ExprChild
		s, operand = parseAtomicExpr(s)

		if left == nil {
			left = operand
		} else {
			left = ast.Binary{left, operator, operand, s.last}
		}

		if multiline {
			for s.skip(token.Newline) {
			}
		}

		var op ast.BinaryOp
		var ok bool
		s, op, ok = parseInfixOperator(s)
		if !ok {
			return s, left
		}

		if secondary && op.Precedence() != operator.Precedence() {
			pan.Panic(newError(s.pos(), "operators have different precedence"))
		}
		operator = op
		secondary = true
	}
}

func parseAtomicExpr(s scan) (scan, ast.ExprChild) {
	return parse(s,
		parseAddress,
		parseCallInExpr,
		parseCharacter,
		parseClone,
		parseFalse,
		parseIndexInExpr,
		parseInteger,
		parseNil,
		parseParenthesized,
		parsePointerDereference,
		parseSelectorInExpr,
		parseString,
		parseTrue,
		parseUnary,
		parseZero,
	)
}

func parseExprList(s scan) (scan, []ast.ExprListChild) {
	if s.skip(token.ParenLeft) {
		return parseListUntil(s, skipper(token.ParenRight),
			parseCommaInExprList,
			parseCommentInExprList,
			parseExpressionInExprList,
			parseNewlineInExprList,
		)
	} else {
		return parseNakedList(s, 0,
			parseCommaInExprList,
			parseExpressionInExprList,
		)
	}
}

func parseAddress(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Ampersand, "address operator expected")
	s, expr := parseAtomicExpr(s)
	return s, ast.Address{t.Pos(), expr, s.last}
}

func parseAssignerDereference(s scan) (scan, ast.AssignerDereference, bool) {
	start := s

	if !s.skip(token.ParenLeft) {
		return start, ast.AssignerDereference{}, false
	}

	t := s.peek()
	if t.Kind != token.Word {
		return start, ast.AssignerDereference{}, false
	}
	s.skip(t.Kind)

	if !s.skip(token.ParenRight) {
		return start, ast.AssignerDereference{}, false
	}

	return s, ast.AssignerDereference{start.pos(), t.Source, s.last}, true
}

func parseAssignerDereferenceInAssignList(s scan) (scan, ast.AssignListChild) {
	s, node, ok := parseAssignerDereference(s)
	if !ok {
		pan.Panic(newError(s.pos(), "assigner dereference expected"))
	}
	return s, node
}

func parseCall(s scan, parsers ...func(scan) (scan, ast.ExprListChild)) (scan, ast.Call) {
	s, name := parseSelectorOnly(s)

	s.take(token.ParenLeft, "call: opening paren expected")
	s, args := parseListUntil(s, skipper(token.ParenRight), parsers...)

	return s, ast.Call{name, args, s.last}
}

func parseCallInExpr(s scan) (scan, ast.ExprChild) {
	return parseCall(s,
		parseCommaInExprList,
		parseCommentInExprList,
		parseExpressionInExprList,
		parseNewlineInExprList,
	)
}

func parseCallInAssignList(s scan) (scan, ast.AssignListChild) {
	return parseCall(s,
		parseCommaInExprList,
		parseExpressionInExprList,
		parseNewlineInExprList,
	)
}

func parseCharacter(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Character, "character literal expected")
	return s, ast.Character{t.Pos(), t.Source, s.last}
}

func parseClone(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Clone, "clone keyword expected")
	s, expr := parseAtomicExpr(s)
	return s, ast.Clone{t.Pos(), expr, s.last}
}

func parseFalse(s scan) (scan, ast.ExprChild) {
	t := s.take(token.False, "literal false expected")
	return s, ast.Boolean{t.Pos(), false, s.last}
}

func parseIndex(s scan) (scan, ast.Index) {
	s, name := parseSelectorOnly(s)
	s.take(token.BracketLeft, "index: opening bracket expected")
	s, index := parseAnyExpr(s, true)
	s.take(token.BracketRight, "index: closing bracket expected")
	return s, ast.Index{name, index, s.last}
}

func parseIndexInAssignList(s scan) (scan, ast.AssignListChild) { return parseIndex(s) }
func parseIndexInExpr(s scan) (scan, ast.ExprChild)             { return parseIndex(s) }

func parseInteger(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Integer, "integer literal expected")
	return s, ast.Integer{t.Pos(), t.Source, s.last}
}

func parseNil(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Nil, "literal nil expected")
	return s, ast.Nil{t.Pos(), s.last}
}

func parseParenthesized(s scan) (scan, ast.ExprChild) {
	s.take(token.ParenLeft, "expression: opening paren expected")
	s, expr := parseAnyExpr(s, true)
	s.take(token.ParenRight, "expression: closing paren expected")
	return s, expr
}

func parsePointerDereference(s scan) (scan, ast.ExprChild) {
	t := s.take(token.Asterisk, "pointer dereference operator expected")
	s, expr := parseAtomicExpr(s)
	return s, ast.PointerDereference{t.Pos(), expr, s.last}
}

func parseSelectorOnly(s scan) (scan, ast.Selector) {
	pos := s.pos()
	name := s.take(token.Word, "selector: variable name expected")

	var names []string

	for {
		names = append(names, name.Source)

		if !s.skip(token.Period) {
			return s, ast.Selector{pos, names, s.last}
		}

		name = s.take(token.Word, "selector: field name expected")
	}
}

func parseSelector(s scan) (scan, ast.Selector) {
	s, name := parseSelectorOnly(s)

	switch s.peek().Kind {
	case token.Colons:
		pan.Panic(newError(s.pos(), "selector: looks like namespace"))
	case token.ParenLeft:
		pan.Panic(newError(s.pos(), "selector used in function call"))
	case token.BracketLeft:
		pan.Panic(newError(s.pos(), "selector: looks like index expression"))
	}

	return s, name
}

func parseSelectorInAssignList(s scan) (scan, ast.AssignListChild) { return parseSelector(s) }
func parseSelectorInExpr(s scan) (scan, ast.ExprChild)             { return parseSelector(s) }

func parseString(s scan) (scan, ast.ExprChild) {
	t := s.take(token.String, "string literal expected")
	return s, ast.String{t.Pos(), t.Source, s.last}
}

func parseTrue(s scan) (scan, ast.ExprChild) {
	t := s.take(token.True, "literal true expected")
	return s, ast.Boolean{t.Pos(), true, s.last}
}

func parseUnary(s scan) (scan, ast.ExprChild) {
	pos := s.pos()
	s, op := parsePrefixOperator(s)
	s, expr := parseAtomicExpr(s)
	return s, ast.Unary{pos, op, expr, s.last}
}

func parseZero(s scan) (scan, ast.ExprChild) {
	t := s.take(token.BraceLeft, "zero: opening brace expected")
	s.take(token.BraceRight, "zero: closing brace expected")
	return s, ast.Zero{t.Pos(), s.last}
}
