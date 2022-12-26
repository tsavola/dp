// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"sort"
	"strings"

	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/internal/namespace"
)

type commentedName struct {
	head []ast.Comment
	name string
	tail []ast.Comment
}

type commentedImport struct {
	head  []ast.Comment
	path  *string
	names []commentedName
	tail  *ast.Comment
}

type commentedListImport struct {
	head  []ast.Comment
	path  string
	names []ast.IdentListChild
	tail  *ast.Comment
}

type commentedImports struct {
	head []ast.Comment
	list []commentedImport
}

type importKey struct {
	path    string
	comment string
}

func makeImportKey(path string, tail []ast.Comment) importKey {
	var comment string
	if len(tail) > 0 {
		comment = strings.TrimSpace(tail[0].Source)
	}
	return importKey{path, comment}
}

func mergeImports(groups []commentedNode[ast.FileChild]) (index int, imports commentedImports) {
	var (
		firstImportIndex    = -1
		firstImportsIndex   = -1
		firstSubstanceIndex = -1

		head []ast.Comment
		list []commentedNode[ast.Import]
	)

	for i, g := range groups {
		if g.node != nil {
			ast.VisitFileChild(*g.node,
				func(ast.Comment) {},

				func(node ast.ConstantDef) {
					if firstSubstanceIndex < 0 {
						firstSubstanceIndex = i
					}
				},

				func(node ast.FunctionDef) {
					if firstSubstanceIndex < 0 {
						firstSubstanceIndex = i
					}

					list = appendImportsFromBlock(list, node.Body)
				},

				func(node ast.Import) {
					if firstImportIndex < 0 {
						firstImportIndex = i
					}

					list = append(list, commentedNode[ast.Import]{g.head, &node, g.tail})
				},

				func(node ast.Imports) {
					if firstImportsIndex < 0 {
						firstImportsIndex = i
					}

					head = append(head, g.head...)
					list = append(list, splitCommentedNodes[ast.ImportListChild, ast.Import](node.Imports, false)...)
				},

				func(node ast.TypeDef) {
					if firstSubstanceIndex < 0 {
						firstSubstanceIndex = i
					}

					for _, node := range node.Fields {
						ast.VisitFieldListChild(node,
							func(ast.Comment) {},
							func(ast.Field) {},
							func(node ast.Import) {
								list = append(list, commentedNode[ast.Import]{nil, &node, nil})
							},
						)
					}
				},
			)
		}
	}

	switch {
	case firstImportsIndex >= 0:
		index = firstImportsIndex
	default:
		index = firstImportIndex
	}

	if firstSubstanceIndex >= 0 && index > firstSubstanceIndex {
		index = firstSubstanceIndex
	}

	return index, commentedImports{head, trimImports(resolveImports(list))}
}

func appendImportsFromBlock(list []commentedNode[ast.Import], nodes []ast.BlockChild) []commentedNode[ast.Import] {
	for _, node := range nodes {
		ast.VisitBlockChild(node,
			func(ast.Assign) {},
			func(node ast.Block) { list = appendImportsFromBlock(list, node.Body) },
			func(ast.Break) {},
			func(ast.Comment) {},
			func(ast.Continue) {},
			func(ast.Expression) {},
			func(ast.For) {},
			func(ast.If) {},
			func(node ast.Import) { list = append(list, commentedNode[ast.Import]{nil, &node, nil}) },
			func(ast.Return) {},
			func(ast.VariableDecl) {},
			func(ast.VariableDef) {},
		)
	}

	return list
}

func resolveImports(groups []commentedNode[ast.Import]) []commentedNode[ast.Import] {
	namespacePaths := make(map[string]*string, len(groups))

	for _, g := range groups {
		if g.node == nil || g.node.Path == "" {
			continue
		}

		path, ok := namespace.UnquoteImportPath(g.node.Path)
		if !ok {
			continue
		}

		for _, s := range namespace.ImportPathNamespaces(path) {
			if value, found := namespacePaths[s]; !found {
				namespacePaths[s] = &g.node.Path
			} else if value != nil && *value != g.node.Path {
				namespacePaths[s] = nil // Disable ambiguous namespace.
			}
		}
	}

	resolved := make([]commentedNode[ast.Import], 0, len(groups))

	for _, g := range groups {
		if g.node == nil || g.node.Path != "" {
			resolved = append(resolved, g)
		} else {
			// TODO: comments
			var (
				pathNames = make(map[string][]ast.IdentListChild)
				badNames  []ast.IdentListChild
			)

			for _, node := range g.node.Names {
				ast.VisitIdentListChild(node,
					func(ast.Comment) {},
					func(node ast.Identifier) {
						if path := namespacePaths[node.Name.Namespace()]; path != nil {
							pathNames[*path] = append(pathNames[*path], ast.Identifier{
								node.Position,
								ast.QualifiedName{node.Name.Short()},
								node.End,
							})
						} else {
							badNames = append(badNames, node)
						}
					},
				)
			}

			for path, names := range pathNames {
				resolved = append(resolved, commentedNode[ast.Import]{
					node: &ast.Import{
						g.node.Position,
						path,
						names,
						g.node.End,
					},
				})
			}

			if len(badNames) > 0 {
				g.node.Names = badNames
				resolved = append(resolved, g)
			}
		}
	}

	return resolved
}

func trimImports(groups []commentedNode[ast.Import]) []commentedImport {
	var (
		merged = make(map[importKey]*commentedListImport, len(groups))
		keys   = make([]importKey, 0, len(groups))
		extra  []ast.Comment
	)

	for _, g := range groups {
		if g.node != nil {
			key := makeImportKey(g.node.Path, g.tail)

			if merg := merged[key]; merg != nil {
				merg.head = append(merg.head, g.head...)
				merg.names = append(merg.names, g.node.Names...)
			} else {
				// Copy ast.Import.Names to avoid mutating FormatFile argument.
				names := append([]ast.IdentListChild(nil), g.node.Names...)

				var tail *ast.Comment
				if len(g.tail) > 0 {
					tail = &g.tail[0]
				}

				merged[key] = &commentedListImport{g.head, g.node.Path, names, tail}
				keys = append(keys, key)
			}
		} else {
			extra = append(extra, g.head...)
		}
	}

	sort.SliceStable(keys, func(i, j int) bool {
		igroup := importPathGroup(keys[i].path)
		jgroup := importPathGroup(keys[j].path)
		if igroup == jgroup {
			return keys[i].path < keys[j].path
		}
		return igroup < jgroup
	})

	list := make([]commentedImport, 0, len(keys)+1)

	for _, key := range keys {
		g := merged[key]
		list = append(list, commentedImport{g.head, &g.path, trimImportNames(g.names), g.tail})
	}

	if len(extra) > 0 {
		list = append(list, commentedImport{extra, nil, nil, nil})
	}

	return list
}

func trimImportNames(nodes []ast.IdentListChild) []commentedName {
	var (
		groups = splitCommentedNodes[ast.IdentListChild, ast.Identifier](nodes, false)
		merged = make(map[string]*commentedName, len(groups))
		names  = make([]string, 0, len(groups))
		extra  []ast.Comment
	)

	for _, g := range groups {
		if g.node != nil {
			name := g.node.Name.String()

			if merg := merged[name]; merg != nil {
				merg.head = append(merg.head, g.head...)
				merg.tail = append(merg.tail, g.tail...)
			} else {
				merged[name] = &commentedName{g.head, name, g.tail}
				names = append(names, name)
			}
		} else {
			extra = append(extra, g.head...)
		}
	}

	sort.SliceStable(names, func(i, j int) bool { return names[i] < names[j] })

	list := make([]commentedName, 0, len(names)+1)

	for _, name := range names {
		list = append(list, *merged[name])
	}

	if len(extra) > 0 {
		list = append(list, commentedName{extra, "", nil})
	}

	return list
}

func importPathGroup(path string) int {
	path, _ = namespace.UnquoteImportPath(path)

	if root, _, _ := strings.Cut(path, "/"); root == "internal" {
		return 3
	}

	if i := strings.Index(path, "."); i >= 0 {
		if !strings.Contains(path[:i], "/") {
			return 2
		}
	}

	return 1
}
