// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"unicode/utf8"

	"github.com/tsavola/dp/ast"
)

func getColumnWidths[T ast.Node](nodes []T, getColumns func(T) []string) map[int][]*int {
	widths := make(map[int][]*int, len(nodes))

	for _, node := range nodes {
		if values := getColumns(node); len(values) > 0 {
			var (
				line = node.Pos().Line
				prev = widths[line-1]
				curr []*int
			)

			for i, value := range values {
				n := len(value) + 1 // Including space.

				var cell *int
				if i < len(prev)-1 {
					cell = prev[i]
					if *cell < n {
						*cell = n
					}
				} else {
					cell = &n
				}

				curr = append(curr, cell)
			}

			widths[line] = curr
		}
	}

	return widths
}

func formatColumns(w writer, values []string, columnWidths []*int) {
	for i := range values {
		if i > 0 {
			n := *columnWidths[i-1]
			pad(w, n-utf8.RuneCountInString(values[i-1]))
		}
		w.WriteString(values[i])
	}
}
