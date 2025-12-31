// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"strings"

	"github.com/tsavola/dp/ast"
)

type commentedNode[T ast.Node] struct {
	head []ast.Comment
	node *T // Non-comment node.
	tail []ast.Comment
}

// splitCommentedNodes file, import list or identifier list children.
func splitCommentedNodes[T ast.Node, R ast.Node](nodes []T, splitOnGap bool) []commentedNode[R] {
	var groups []commentedNode[R]
	var g commentedNode[R]

	for i := 0; i < len(nodes); i++ {
		curr := nodes[i]

		switch x := ast.Node(curr).(type) {
		case ast.Comment:
			g.head = append(g.head, x)
		case R:
			g.node = &x
		default:
			panic(curr)
		}

		if i+1 == len(nodes) {
			break
		}

		next := nodes[i+1]
		step := next.Pos().Line - curr.EndPos().Line

		var split bool

		switch {
		case splitOnGap && step >= 2:
			split = true

		case g.node != nil:
			if step == 0 {
				if c, ok := ast.Node(next).(ast.Comment); ok {
					switch ast.Node(curr).(type) {
					case ast.ConstantDef:
						// Preserve comment on same line.
						// TODO: preserve also on following lines?
						g.tail = []ast.Comment{c}
						i++

					case ast.Identifier:
						// Preserve comments on same and following lines.
						for ok {
							g.tail = append(g.tail, c)
							i++

							if i+1 == len(nodes) {
								break
							}

							curr = next
							next = nodes[i+1]
							step = next.Pos().Line - curr.EndPos().Line

							if step != 1 {
								break
							}

							c, ok = ast.Node(next).(ast.Comment)
						}

					case ast.Import:
						// Preserve comment on same line.
						g.tail = []ast.Comment{c}
						i++
					}
				}
			}

			split = true
		}

		if split {
			groups = append(groups, g)
			g = commentedNode[R]{}
		}
	}

	if len(g.head) > 0 || g.node != nil {
		groups = append(groups, g)
	}

	return groups
}

// formatComment grows offsets if necessary.
func formatComment(w writer, level int, node ast.Comment, offsets map[int]*int) {
	lineLen := w.currentLineLen()

	if off := offsets[node.Line]; off != nil {
		pad(w, *off-lineLen)
	} else {
		off = offsets[node.Line-1]

		switch {
		case off == nil: // Previous line has no comment.
			off = &lineLen

		case *off <= level: // Previous line has only a comment.
			off = &lineLen

		case *off < lineLen:
			*off = lineLen
		}

		offsets[node.Line] = off
	}

	formatCommentAlone(w, node)
}

func formatCommentAlone(w writer, node ast.Comment) {
	w.WriteString(strings.TrimSpace(node.Source))
}
