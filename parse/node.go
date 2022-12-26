// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/token"
)

func parseComma(s scan) scan {
	s.take(token.Comma, "comma expected")
	return s
}

func parseCommaInAssignList(s scan) (scan, ast.AssignListChild) { return parseComma(s), nil }
func parseCommaInExprList(s scan) (scan, ast.ExprListChild)     { return parseComma(s), nil }
func parseCommaInFieldList(s scan) (scan, ast.FieldListChild)   { return parseComma(s), nil }
func parseCommaInIdentList(s scan) (scan, ast.IdentListChild)   { return parseComma(s), nil }
func parseCommaInImportList(s scan) (scan, ast.ImportListChild) { return parseComma(s), nil }
func parseCommaInParamList(s scan) (scan, ast.ParamListChild)   { return parseComma(s), nil }
func parseCommaInTypeList(s scan) (scan, ast.TypeListChild)     { return parseComma(s), nil }
func parseCommaString(s scan) (scan, string)                    { return parseComma(s), "" }

func parseComment(s scan) (scan, ast.Comment) {
	t := s.take(token.Comment, "comment expected")
	return s, ast.Comment{t.Pos(), t.Source}
}

func parseCommentInBlock(s scan) (scan, ast.BlockChild)           { return parseComment(s) }
func parseCommentInExprList(s scan) (scan, ast.ExprListChild)     { return parseComment(s) }
func parseCommentInFieldList(s scan) (scan, ast.FieldListChild)   { return parseComment(s) }
func parseCommentInFile(s scan) (scan, ast.FileChild)             { return parseComment(s) }
func parseCommentInIdentList(s scan) (scan, ast.IdentListChild)   { return parseComment(s) }
func parseCommentInImportList(s scan) (scan, ast.ImportListChild) { return parseComment(s) }
func parseCommentInParamList(s scan) (scan, ast.ParamListChild)   { return parseComment(s) }
func parseCommentInTypeList(s scan) (scan, ast.TypeListChild)     { return parseComment(s) }

func parseNewline(s scan) scan {
	s.take(token.Newline, "end of line expected")
	return s
}

func parseNewlineInBlock(s scan) (scan, ast.BlockChild)           { return parseNewline(s), nil }
func parseNewlineInExprList(s scan) (scan, ast.ExprListChild)     { return parseNewline(s), nil }
func parseNewlineInFieldList(s scan) (scan, ast.FieldListChild)   { return parseNewline(s), nil }
func parseNewlineInFile(s scan) (scan, ast.FileChild)             { return parseNewline(s), nil }
func parseNewlineInIdentList(s scan) (scan, ast.IdentListChild)   { return parseNewline(s), nil }
func parseNewlineInImportList(s scan) (scan, ast.ImportListChild) { return parseNewline(s), nil }
func parseNewlineInParamList(s scan) (scan, ast.ParamListChild)   { return parseNewline(s), nil }
func parseNewlineInTypeList(s scan) (scan, ast.TypeListChild)     { return parseNewline(s), nil }

func parseSemicolon(s scan) scan {
	s.take(token.Semicolon, "semicolon expected")
	return s
}

func parseSemicolonInBlock(s scan) (scan, ast.BlockChild)           { return parseSemicolon(s), nil }
func parseSemicolonInFieldList(s scan) (scan, ast.FieldListChild)   { return parseSemicolon(s), nil }
func parseSemicolonInFile(s scan) (scan, ast.FileChild)             { return parseSemicolon(s), nil }
func parseSemicolonInImportList(s scan) (scan, ast.ImportListChild) { return parseSemicolon(s), nil }
