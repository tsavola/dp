// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
	"import.name/pan"
)

func parseStatements(s scan) (scan, []ast.BlockChild) {
	return parseListUntil(s, skipper(token.BraceRight),
		parseAssign,
		parseBlock,
		parseBreak,
		parseCommentInBlock,
		parseContinue,
		parseExpressionInBlock,
		parseFor,
		parseIf,
		parseImportInBlock,
		parseNewlineInBlock,
		parseReturn,
		parseSemicolonInBlock,
		parseVariableDecl,
		parseVariableDef,
	)
}

func parseAssign(s scan) (scan, ast.BlockChild) {
	s, names := parseNakedList(s, token.Assign,
		parseAssignerDereferenceInAssignList,
		parseCallInAssignList,
		parseCommaInAssignList,
		parseIndexInAssignList,
		parseSelectorInAssignList,
	)
	if len(names) == 0 {
		pan.Panic(newError(s.pos(), "assign: empty list"))
	}

	s.take(token.Assign, "assign: operator expected")

	s, values := parseExprList(s)
	if len(values) == 0 {
		pan.Panic(newError(s.pos(), "assign: empty list"))
	}

	return s, ast.Assign{names[0].Pos(), names, values, s.last}
}

func parseBlock(s scan) (scan, ast.BlockChild) {
	t := s.take(token.BraceLeft, "block: opening brace expected")
	s, body := parseStatements(s)
	return s, ast.Block{t.Pos(), body, s.last}
}

func parseBreak(s scan) (scan, ast.BlockChild) {
	t := s.take(token.Break, "break keyword expected")
	return s, ast.Break{t.Pos(), s.last}
}

func parseContinue(s scan) (scan, ast.BlockChild) {
	t := s.take(token.Continue, "continue keyword expected")
	return s, ast.Continue{t.Pos(), s.last}
}

func parseFor(s scan) (scan, ast.BlockChild) {
	keyword := s.take(token.For, "for keyword expected")

	var test ast.ExprChild

	open, ok := s.skim(token.BraceLeft)
	if !ok {
		s, test = parseAnyExpr(s, false)
		open = s.take(token.BraceLeft, "for: opening brace expected")
	}

	s, body := parseStatements(s)
	return s, ast.For{keyword.Pos(), test, open.Pos(), body, s.last}
}

func parseIf(s scan) (scan, ast.BlockChild) {
	keyword := s.take(token.If, "if keyword expected")
	s, test := parseAnyExpr(s, false)

	thenPos := s.take(token.BraceLeft, "if: opening brace expected").Pos()
	s, then := parseStatements(s)
	thenEnd := s.last

	var els []ast.BlockChild
	if s.skip(token.Else) {
		s.take(token.BraceLeft, "else: opening brace expected")
		s, els = parseStatements(s)
	}

	return s, ast.If{keyword.Pos(), test, thenPos, then, thenEnd, els, s.last}
}

func parseReturn(s scan) (scan, ast.BlockChild) {
	keyword := s.take(token.Return, "return keyword expected")

	s, values := parse(s,
		func(s scan) (scan, []ast.ExprListChild) {
			return parseNakedList(s, 0,
				parseCommaInExprList,
				parseExpressionInExprList,
			)
		},

		func(s scan) (scan, []ast.ExprListChild) {
			s.take(token.ParenLeft, "return value list expected")
			return parseListUntil(s, skipper(token.ParenRight),
				parseCommaInExprList,
				parseCommentInExprList,
				parseExpressionInExprList,
				parseNewlineInExprList,
			)
		},
	)

	return s, ast.Return{keyword.Pos(), values, s.last}
}

func parseVariableDecl(s scan) (scan, ast.BlockChild) {
	pos := s.pos()

	s, names := parseNakedList(s, 0,
		parseCommaString,
		parseVariableName,
	)

	s.take(token.Colon, "variable declaration: colon expected")

	if s.skip(token.Auto) {
		return s, ast.VariableDecl{pos, names, nil, s.last}
	}

	s, spec := parseTypeSpec(s)

	return s, ast.VariableDecl{pos, names, &spec, s.last}
}

func parseVariableDef(s scan) (scan, ast.BlockChild) {
	pos := s.pos()

	s, names := parseNakedList(s, 0,
		parseCommaString,
		parseVariableName,
	)

	s.take(token.Define, "variable definition: operator expected")

	s, values := parseExprList(s)

	return s, ast.VariableDef{pos, names, values, s.last}
}

func parseVariableName(s scan) (scan, string) {
	name := s.take(token.Word, "variable name expected")
	return s, name.Source
}
