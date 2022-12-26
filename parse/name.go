// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
)

func parseQualifiedName(s scan) (scan, ast.QualifiedName) {
	var name ast.QualifiedName

	if s.skip(token.Colons) {
		name = append(name, "")
	}

	for {
		name = append(name, s.take(token.Word, "name expected").Source)

		if !s.skip(token.Colons) {
			return s, name
		}
	}
}

func parseIdentifierInIdentList(s scan) (scan, ast.IdentListChild) {
	pos := s.pos()
	s, name := parseQualifiedName(s)
	return s, ast.Identifier{pos, name, s.last}
}
