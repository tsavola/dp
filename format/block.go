// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"strings"

	"github.com/tsavola/dp/ast"
)

func formatBlock(w writer, level, startLine int, nodes []ast.BlockChild) {
	w.WriteString("{")

	columnify := func(node ast.BlockChild) []string {
		switch node := node.(type) {
		case ast.VariableDecl:
			return []string{strings.Join(node.Names, ", "), ":"}
		case ast.VariableDef:
			return []string{strings.Join(node.Names, ", "), ":="}
		default:
			return nil
		}
	}

	var (
		columnWidths   = getColumnWidths(nodes, columnify)
		commentOffsets = make(map[int]*int)
	)

	base := w.Len()
	formatStatements(w, level, startLine, nodes, columnify, columnWidths, commentOffsets)
	w.Truncate(base)
	formatStatements(w, level, startLine, nodes, columnify, columnWidths, commentOffsets)

	w.WriteString("\n")
	indent(w, level-1)
	w.WriteString("}")
}

func formatStatements(
	w writer,
	level int,
	startLine int,
	nodes []ast.BlockChild,
	columnify func(ast.BlockChild) []string,
	columnWidths map[int][]*int,
	commentOffsets map[int]*int,
) {
	prevLine := startLine

	for i, node := range nodes {
		indentNode(w, level, prevLine, node)

		ast.VisitBlockChild(node,
			func(node ast.Assign) {
				for i, node := range node.Objects {
					if i > 0 {
						w.WriteString(", ")
					}
					formatAssignListChild(w, node)
				}
				w.WriteString(" = ")
				formatExprList(w, level+1, node.Line, node.Subjects, false)
			},

			func(node ast.Block) {
				formatBlock(w, level+1, node.Pos().Line, node.Body)
			},

			func(node ast.Break) {
				w.WriteString("break")
			},

			func(node ast.Comment) {
				formatComment(w, level, node, i, commentOffsets)
			},

			func(node ast.Continue) {
				w.WriteString("continue")
			},

			func(node ast.Expression) {
				formatExpr(w, level+1, node.Expr, 0, false)
			},

			func(node ast.For) {
				w.WriteString("for ")
				if node.Test != nil {
					formatExpr(w, level+1, node.Test, 0, false)
					w.WriteString(" ")
				}
				formatBlock(w, level+1, node.BodyPos.Line, node.Body)
			},

			func(node ast.If) {
				w.WriteString("if ")
				formatExpr(w, level+1, node.Test, 0, false)
				w.WriteString(" ")
				formatBlock(w, level+1, node.ThenPos.Line, node.Then)
				if len(node.Else) > 0 {
					w.WriteString(" else ")
					formatBlock(w, level+1, node.ThenEnd.Line, node.Else)
				}
			},

			func(ast.Import) {},

			func(node ast.Return) {
				w.WriteString("return")
				if len(node.Values) > 0 {
					w.WriteString(" ")
				}
				formatExprList(w, level+1, node.Line, node.Values, false)
			},

			func(node ast.VariableDecl) {
				formatColumns(w, columnify(node), columnWidths[node.Line])
				w.WriteString(" ")
				if node.Type == nil {
					w.WriteString("auto")
				} else {
					w.WriteString(node.Type.Type.String())
				}
			},

			func(node ast.VariableDef) {
				formatColumns(w, columnify(node), columnWidths[node.Line])
				w.WriteString(" ")
				formatExprList(w, level+1, node.Line, node.Values, false)
			},
		)

		if node.Pos().Line != node.EndPos().Line {
			// Discontinue comment lineage after multi-line statement.
			commentOffsets[node.EndPos().Line] = new(int)
		}

		prevLine = node.EndPos().Line
	}
}

func formatAssignListChild(w writer, node ast.AssignListChild) {
	ast.VisitAssignListChild(node,
		func(node ast.AssignerDereference) {
			formatAssignerDereference(w, node)
		},
		func(node ast.Call) {
			w.WriteString(node.Name.String())
			formatExprList(w, 0, node.Pos().Line, node.Args, true)
		},
		func(node ast.Index) {
			w.WriteString(node.Name.String())
			w.WriteString("[")
			formatExpr(w, 0, node.Index, 0, true)
			w.WriteString("]")
		},
		func(node ast.Selector) {
			w.WriteString(node.String())
		},
	)
}
