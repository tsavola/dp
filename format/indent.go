// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"github.com/tsavola/dp/ast"
)

func pad(w writer, amount int) {
	for range amount {
		w.WriteString(" ")
	}
}

func indent(w writer, level int) {
	for range level {
		w.WriteString("\t")
	}
}

func indentNode(w writer, level, prevLine int, node ast.Node) {
	step := node.Pos().Line - prevLine

	switch {
	case step == 0 && ast.IsComment(node):
		w.WriteString(" ")
	case step <= 1:
		w.WriteString("\n")
		indent(w, level)
	default:
		w.WriteString("\n\n")
		indent(w, level)
	}
}
