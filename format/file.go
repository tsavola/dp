// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package format generates canonically formatted source code.
package format

import (
	"bytes"
	"strings"

	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/field"
)

func File(nodes []ast.FileChild) []byte {
	if len(nodes) == 0 {
		return nil
	}

	size := nodes[len(nodes)-1].EndPos().ByteOffset
	w := writer{bytes.NewBuffer(make([]byte, 0, size+size/4))}

	groups := splitCommentedNodes[ast.FileChild, ast.FileChild](nodes, true)
	importsIndex, imports := mergeImports(groups)

	for i, g := range groups {
		if i != importsIndex && g.node != nil {
			var isImport bool

			switch (*g.node).(type) {
			case ast.Import:
				isImport = true
			case ast.Imports:
				isImport = true
			}

			if isImport {
				continue
			}
		}

		if i > 0 {
			gap := true

			// Special case: no forced gap between value with same visibility.
			prev := groups[i-1].node
			curr := g.node
			if curr != nil && prev != nil {
				if prev, ok := (*prev).(ast.ConstantDef); ok {
					if curr, ok := (*curr).(ast.ConstantDef); ok {
						if curr.Line-prev.End.Line <= 1 && curr.Public == prev.Public {
							gap = false
						}
					}
				}
			}

			if gap {
				w.WriteString("\n")
			}
		}

		if i == importsIndex {
			for _, node := range imports.head {
				w.WriteString(strings.TrimSpace(node.Source))
				w.WriteString("\n")
			}

			w.WriteString("import {\n")

			for i, imp := range imports.list {
				if imports.list[i].path == nil || (i > 0 && importPathGroup(*imports.list[i-1].path) != importPathGroup(*imports.list[i].path)) {
					w.WriteString("\n")
				}

				for _, node := range imp.head {
					w.WriteString("\t")
					w.WriteString(strings.TrimSpace(node.Source))
					w.WriteString("\n")
				}

				if imp.path != nil {
					w.WriteString("\t")
					w.WriteString(*imp.path)

					if len(imp.names) > 0 {
						w.WriteString(" (")

						first := true
						for _, g := range imp.names {
							// TODO: comments, multiple lines
							if g.name != "" {
								if !first {
									w.WriteString(", ")
								}
								w.WriteString(g.name)
								first = false
							}
						}

						w.WriteString(")")
					}

					if imp.tail != nil {
						w.WriteString(" ")
						w.WriteString(strings.TrimSpace(imp.tail.Source))
					}
					w.WriteString("\n")
				}
			}

			w.WriteString("}\n")
		}

		for _, node := range g.head {
			w.WriteString(strings.TrimSpace(node.Source))
			w.WriteString("\n")
		}

		if g.node != nil {
			ast.VisitFileChild(*g.node,
				func(ast.Comment) {},

				func(node ast.ConstantDef) {
					if node.Public {
						w.WriteString("pub ")
					}
					w.WriteString(node.ConstName)
					w.WriteString(" = ")
					formatExpr(w, 1, node.Value, 0, false)
					w.WriteString("\n")
				},

				func(node ast.FunctionDef) {
					if node.Public {
						w.WriteString("pub ")
					}
					if node.ReceiverType != nil {
						w.WriteString("(")
						if node.ReceiverName != "" {
							w.WriteString(node.ReceiverName)
							w.WriteString(" ")
						}
						w.WriteString(node.ReceiverType.Type.String())
						w.WriteString(") ")
					}
					w.WriteString(node.FuncName)
					formatFunctionParams(w, node)
					formatFunctionResults(w, node)
					if body := trimFunctionBody(node); len(body) == 0 {
						w.WriteString("{}")
					} else {
						formatBlock(w, 1, node.BodyPos.Line, body)
					}
					w.WriteString("\n")
				},

				func(ast.Import) {},
				func(ast.Imports) {},

				func(node ast.TypeDef) {
					empty := true
					for _, node := range node.Fields {
						ast.VisitFieldListChild(node,
							func(node ast.Comment) { empty = false },
							func(node ast.Field) { empty = false },
							func(ast.Import) {},
						)
					}
					formatTypeDefBody(w, node, empty)
					w.WriteString("\n")
				},
			)
		}
	}

	return w.Bytes()
}

func formatFunctionParams(w writer, def ast.FunctionDef) {
	var (
		comments bool
		params   []ast.Parameter
	)
	for _, node := range def.Params {
		ast.VisitParamListChild(node,
			func(node ast.Comment) { comments = true },
			func(node ast.Parameter) { params = append(params, node) },
		)
	}

	w.WriteString("(")

	if !comments && (len(params) == 0 || params[0].Line == def.Line) {
		for i, param := range params {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(param.ParamName)
			if i == len(params)-1 || !param.Type.Type.Equal(params[i+1].Type.Type) {
				w.WriteString(" ")
				w.WriteString(param.Type.String())
			}
		}
	} else {
		columnify := func(node ast.ParamListChild) []string {
			if node, ok := node.(ast.Parameter); ok {
				return []string{node.ParamName, node.Type.Type.String()}
			}
			return nil
		}

		var (
			columnWidths   = getColumnWidths(def.Params, columnify)
			commentOffsets = make(map[int]*int)
		)

		base := w.Len()
		formatFunctionParamsMultiLine(w, def, columnify, columnWidths, commentOffsets)
		w.Truncate(base)
		formatFunctionParamsMultiLine(w, def, columnify, columnWidths, commentOffsets)

		w.WriteString("\n")
	}

	w.WriteString(") ")
}

func formatFunctionParamsMultiLine(w writer, def ast.FunctionDef, columnify func(ast.ParamListChild) []string, columnWidths map[int][]*int, commentOffsets map[int]*int) {
	prevLine := def.Line

	for i, node := range def.Params {
		indentNode(w, 1, prevLine, node)

		ast.VisitParamListChild(node,
			func(node ast.Comment) {
				formatComment(w, 1, node, i, commentOffsets)
			},

			func(node ast.Parameter) {
				formatColumns(w, columnify(node), columnWidths[node.Line])
				w.WriteString(",")
			},
		)

		prevLine = node.EndPos().Line
	}
}

func formatFunctionResults(w writer, def ast.FunctionDef) {
	var (
		comments bool
		specs    []ast.TypeSpec
	)
	for _, node := range def.Results {
		ast.VisitTypeListChild(node,
			func(node ast.Comment) { comments = true },
			func(node ast.TypeSpec) { specs = append(specs, node) },
		)
	}

	switch {
	case !comments && len(specs) == 0:
		w.WriteString("() ")

	case !comments && len(specs) == 1:
		w.WriteString(specs[0].Type.String())
		w.WriteString(" ")

	case !comments && specs[0].Line == def.ParamsEnd.Line:
		w.WriteString("(")

		for i, spec := range specs {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(spec.Type.String())
		}

		w.WriteString(") ")

	default:
		w.WriteString("(")

		commentOffsets := make(map[int]*int)

		base := w.Len()
		formatFunctionResultsMultiLine(w, def, commentOffsets)
		w.Truncate(base)
		formatFunctionResultsMultiLine(w, def, commentOffsets)

		w.WriteString("\n) ")
	}
}

func formatFunctionResultsMultiLine(w writer, def ast.FunctionDef, commentOffsets map[int]*int) {
	prevLine := def.ParamsEnd.Line

	for i, node := range def.Results {
		indentNode(w, 1, prevLine, node)

		ast.VisitTypeListChild(node,
			func(node ast.Comment) {
				formatComment(w, 1, node, i, commentOffsets)
			},

			func(node ast.TypeSpec) {
				w.WriteString(node.Type.String())
				w.WriteString(",")
			},
		)

		prevLine = node.EndPos().Line
	}
}

// trimFunctionBody removes unnecessary return statements.
func trimFunctionBody(def ast.FunctionDef) []ast.BlockChild {
	if len(def.Results) > 0 {
		return def.Body
	}

	nodes := append([]ast.BlockChild{}, def.Body...)

	for i := len(nodes) - 1; i >= 0; i-- {
		switch node := nodes[i].(type) {
		case ast.Comment:

		case ast.Return:
			if len(node.Values) > 0 {
				return nodes
			}

			nodes = append(nodes[:i], nodes[i+1:]...)

		default:
			return nodes
		}
	}

	return nodes
}

func formatTypeDefBody(w writer, node ast.TypeDef, empty bool) {
	if node.Public {
		w.WriteString("pub ")
	}
	w.WriteString(node.TypeName)
	w.WriteString(" {")

	if !empty {
		columnify := func(node ast.FieldListChild) []string {
			if node, ok := node.(ast.Field); ok {
				values := make([]string, 0, 3)
				values = append(values, node.FieldName)
				values = append(values, node.Type.Type.String())
				if node.Access != field.AccessHidden {
					values = append(values, node.Access.String())
				}
				return values
			}
			return nil
		}

		var (
			columnWidths   = getColumnWidths(node.Fields, columnify)
			commentOffsets = make(map[int]*int)
		)

		base := w.Len()
		formatTypeFields(w, node, columnify, columnWidths, commentOffsets)
		w.Truncate(base)
		formatTypeFields(w, node, columnify, columnWidths, commentOffsets)

		w.WriteString("\n")
	}

	w.WriteString("}")
}

func formatTypeFields(w writer, def ast.TypeDef, columnify func(ast.FieldListChild) []string, columnWidths map[int][]*int, commentOffsets map[int]*int) {
	prevLine := def.Line

	for i, node := range def.Fields {
		indentNode(w, 1, prevLine, node)

		ast.VisitFieldListChild(node,
			func(node ast.Comment) { formatComment(w, 1, node, i, commentOffsets) },
			func(node ast.Field) { formatColumns(w, columnify(node), columnWidths[node.Line]) },
			func(ast.Import) {},
		)

		prevLine = node.EndPos().Line
	}
}
