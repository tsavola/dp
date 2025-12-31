// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"github.com/tsavola/dp/ast"
)

func useMultipleLines[T ast.Node](startLine int, nodes ...T) bool {
	for _, node := range nodes {
		if node.EndPos().Line > startLine || ast.IsComment(node) {
			return true
		}
	}
	return false
}

func formatExprList(w writer, level, startLine int, nodes []ast.ExprListChild, forceParens bool) {
	if useMultipleLines(startLine, nodes...) {
		w.WriteString("(")

		commentOffsets := make(map[int]*int)

		base := w.Len()
		formatExprListMultiLine(w, level, startLine, nodes, commentOffsets)
		w.Truncate(base)
		formatExprListMultiLine(w, level, startLine, nodes, commentOffsets)

		w.WriteString("\n")
		indent(w, level-1)
		w.WriteString(")")
	} else {
		if forceParens {
			w.WriteString("(")
		}

		formatExprListOneLine(w, nodes)

		if forceParens {
			w.WriteString(")")
		}
	}
}

func formatExprListOneLine(w writer, nodes []ast.ExprListChild) {
	for i, node := range nodes {
		if i > 0 {
			w.WriteString(", ")
		}

		ast.VisitExprListChild(node,
			func(node ast.AssignerDereference) { formatAssignerDereference(w, node) },
			func(ast.Comment) {},
			func(node ast.Expression) { formatExpr(w, 0, node.Expr, 0, false) },
		)
	}
}

func formatExprListMultiLine(w writer, level, startLine int, nodes []ast.ExprListChild, commentOffsets map[int]*int) {
	first := true
	prevLine := startLine

	for _, node := range nodes {
		indentNode(w, level, prevLine, node)

		ast.VisitExprListChild(node,
			func(node ast.AssignerDereference) {
				formatAssignerDereference(w, node)
				w.WriteString(",")
			},

			func(node ast.Comment) {
				if first {
					formatCommentAlone(w, node)
				} else {
					formatComment(w, 1, node, commentOffsets)
				}
			},

			func(node ast.Expression) {
				formatExpr(w, level+1, node.Expr, 0, false)
				w.WriteString(",")
			},
		)

		first = false
		prevLine = node.EndPos().Line
	}
}
