// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
)

func parseType(s scan) (scan, ast.Type) {
	var t ast.Type

	t.Assigner = s.skip(token.Assign)
	t.Pointer = s.skip(token.Asterisk)
	t.Reference = s.skip(token.Ampersand)
	t.Shared = s.skip(token.Hash)

	if s.skip(token.BracketLeft) {
		var item ast.Type
		s, item = parseType(s)
		t.Item = &item
		s.take(token.BracketRight, "type: array closing bracket expected")
	} else {
		s, t.Name = parseQualifiedName(s)
	}

	return s, t
}

func parseTypeSpec(s scan) (scan, ast.TypeSpec) {
	pos := s.pos()
	s, t := parseType(s)
	return s, ast.TypeSpec{pos, t, s.last}
}

func parseTypeSpecInTypeList(s scan) (scan, ast.TypeListChild) { return parseTypeSpec(s) }

func typeSpecified(t ast.TypeSpec) bool {
	return t.Item != nil || len(t.Name) > 0
}
