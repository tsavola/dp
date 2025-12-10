// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package parse implements parser.
package parse

import (
	"strings"

	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/field"
	"github.com/tsavola/dp/internal/pan"
	"github.com/tsavola/dp/source"
	"github.com/tsavola/dp/token"
)

func File(tokens []token.Token) ([]ast.FileChild, error) {
	var nodes []ast.FileChild

	err := pan.Recover(func() {
		nodes = parseTokens(tokens)
	})

	return nodes, err
}

func parseTokens(tokens []token.Token) []ast.FileChild {
	_, nodes := parseListUntil(scan{tokens, source.Position{}}, peekEOF,
		parseCommentInFile,
		parseConstantDef,
		parseFunctionDef,
		parseImports,
		parseNewlineInFile,
		parseSemicolonInFile,
		parseTypeDef,
	)
	return nodes
}

func parseConstantDef(s scan) (scan, ast.FileChild) {
	var public bool
	var name token.Token

	t := s.take(token.Word, "constant definition: pub keyword or name expected")
	if t.Source == "pub" {
		public = true
		name = s.take(token.Word, "constant definition: name expected")
	} else {
		name = t
	}

	s.take(token.Assign, "constant definition: assignment operator expected")
	s, value := parseAnyExpr(s, false)

	return s, ast.ConstantDef{t.Pos(), public, name.Source, value, s.last}
}

func parseFieldAccess(s scan) (scan, field.Access) {
	if t, ok := s.skim(token.Word); ok {
		switch t.Source {
		case "visible":
			return s, field.AccessVisible
		case "mutable":
			return s, field.AccessMutable
		case "assignable":
			return s, field.AccessAssignable
		}
		pan.Panic(newError(t.Pos(), "visible, mutable or assignable keyword expected"))
	}
	return s, field.AccessHidden
}

func parseFieldInFieldList(s scan) (scan, ast.FieldListChild) {
	name := s.take(token.Word, "field name expected")
	s, spec := parseTypeSpec(s)
	s, access := parseFieldAccess(s)
	return s, ast.Field{name.Pos(), name.Source, spec, access, s.last}
}

func parseFunctionDef(s scan) (scan, ast.FileChild) {
	pos := s.pos()

	var (
		public bool
		rName  string
		rType  *ast.TypeSpec
		name   string
	)

	if t := s.peek(); t.Kind == token.Word && t.Source == "pub" {
		s.skip(token.Word)
		public = true
	}

	if s.skip(token.ParenLeft) {
		rName = s.take(token.Word, "function definition: receiver name expected").Source

		var t ast.TypeSpec
		s, t = parseTypeSpec(s)
		rType = &t

		s.take(token.ParenRight, "function definition: receiver: closing paren expected")

		if t, ok := s.skim(token.Word); ok {
			name = t.Source
		}
	} else {
		name = s.take(token.Word, "function definition: name expected").Source
	}

	s.take(token.ParenLeft, "function definition: parameter list expected")
	s, params := parseListUntil(s, skipper(token.ParenRight),
		parseCommaInParamList,
		parseCommentInParamList,
		parseNewlineInParamList,
		parseParamInParamList,
	)
	paramsEnd := s.last

	params = fillInFunctionParamTypes(params)

	s, results := parse(s,
		func(s scan) (scan, []ast.TypeListChild) {
			return parseNakedList(s, token.BraceLeft,
				parseCommaInTypeList,
				parseTypeSpecInTypeList,
			)
		},

		func(s scan) (scan, []ast.TypeListChild) {
			s.take(token.ParenLeft, "function definition: return type list expected")
			return parseListUntil(s, skipper(token.ParenRight),
				parseCommaInTypeList,
				parseCommentInTypeList,
				parseNewlineInTypeList,
				parseTypeSpecInTypeList,
			)
		},
	)

	bodyPos := s.take(token.BraceLeft, "function definition: opening brace expected").Pos()
	s, body := parseStatements(s)

	return s, ast.FunctionDef{pos, public, rName, rType, name, params, paramsEnd, results, bodyPos, body, s.last}
}

func fillInFunctionParamTypes(nodes []ast.ParamListChild) []ast.ParamListChild {
	var latest ast.TypeSpec

	for i := len(nodes) - 1; i >= 0; i-- {
		ast.VisitParamListChild(nodes[i],
			func(ast.Comment) {},
			func(node ast.Parameter) {
				if typeSpecified(node.Type) {
					latest = node.Type
				} else {
					if !typeSpecified(latest) {
						pan.Panic(newError(node.EndPos(), "function parameter type expected"))
					}
					node.Type = latest
					nodes[i] = node
				}
			},
		)
	}

	return nodes
}

func parseImport(s scan, requireKeyword bool) (scan, ast.Import) {
	pos := s.pos()

	if !s.skip(token.Import) && requireKeyword {
		pan.Panic(newError(pos, "import keyword expected"))
	}

	var path string
	if t, ok := s.skim(token.String); ok {
		path = t.Source
		if !strings.HasPrefix(path, `"`) {
			pan.Panic(newError(t.Pos(), "import path: opening quote expected"))
		}
	}

	var names []ast.IdentListChild
	if s.skip(token.ParenLeft) {
		s, names = parseListUntil(s, skipper(token.ParenRight),
			parseCommaInIdentList,
			parseCommentInIdentList,
			parseIdentifierInIdentList,
			parseNewlineInIdentList,
		)
	} else {
		s, names = parseNakedList(s, 0,
			parseCommaInIdentList,
			parseIdentifierInIdentList,
		)
	}

	if !requireKeyword && path == "" && len(names) == 0 {
		pan.Panic(newError(s.pos(), "import: path or identifier list expected"))
	}

	return s, ast.Import{pos, path, names, s.last}
}

func parseImportInBlock(s scan) (scan, ast.BlockChild)           { return parseImport(s, true) }
func parseImportInFieldList(s scan) (scan, ast.FieldListChild)   { return parseImport(s, true) }
func parseImportInImportList(s scan) (scan, ast.ImportListChild) { return parseImport(s, false) }

func parseImportPathInImportList(s scan) (scan, ast.ImportListChild) {
	path := s.take(token.String, "import path expected")
	return s, ast.Import{path.Pos(), path.Source, nil, s.last}
}

func parseImports(s scan) (scan, ast.FileChild) {
	t := s.take(token.Import, "import keyword expected")

	var imports []ast.ImportListChild

	switch {
	case s.peek().Kind == token.String:
		s, imports = parseNakedList(s, 0,
			parseCommaInImportList,
			parseImportPathInImportList,
		)

	case s.skip(token.ParenLeft):
		s, imports = parseListUntil(s, skipper(token.ParenRight),
			parseCommaInImportList,
			parseCommentInImportList,
			parseImportPathInImportList,
			parseNewlineInImportList,
		)

	default:
		s.take(token.BraceLeft, "import: opening brace expected")
		s, imports = parseListUntil(s, skipper(token.BraceRight),
			parseCommaInImportList,
			parseCommentInImportList,
			parseImportInImportList,
			parseNewlineInImportList,
			parseSemicolonInImportList,
		)
	}

	return s, ast.Imports{t.Pos(), imports, s.last}
}

func parseParamInParamList(s scan) (scan, ast.ParamListChild) {
	name := s.take(token.Word, "parameter name expected")

	// Missing type is filled in by parseFunctionDef().
	var spec ast.TypeSpec
	if !s.skip(token.Comma) {
		s, spec = parseTypeSpec(s)
		s.skim(token.Comma)
	}

	return s, ast.Parameter{name.Pos(), name.Source, spec, s.last}
}

func parseTypeDef(s scan) (scan, ast.FileChild) {
	var public bool
	var name token.Token

	t := s.take(token.Word, "type definition: pub keyword or name expected")
	if t.Source == "pub" {
		public = true
		name = s.take(token.Word, "type definition: name expected")
	} else {
		name = t
	}

	s, access := parseFieldAccess(s)

	s.take(token.BraceLeft, "type definition: opening brace expected")
	s, body := parseListUntil(s, skipper(token.BraceRight),
		parseCommaInFieldList,
		parseCommentInFieldList,
		parseFieldInFieldList,
		parseImportInFieldList,
		parseNewlineInFieldList,
		parseSemicolonInFieldList,
	)

	body = fillInTypeFieldAccess(body, access)

	return s, ast.TypeDef{t.Pos(), public, name.Source, body, s.last}
}

func fillInTypeFieldAccess(nodes []ast.FieldListChild, fallback field.Access) []ast.FieldListChild {
	if fallback != 0 {
		for i, node := range nodes {
			ast.VisitFieldListChild(node,
				func(ast.Comment) {},
				func(node ast.Field) {
					if node.Access == 0 {
						node.Access = fallback
						nodes[i] = node
					}
				},
				func(ast.Import) {},
			)
		}
	}

	return nodes
}
