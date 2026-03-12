// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"strings"

	"github.com/tsavola/dp/field"
	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
)

type Comment struct {
	At     source.Position
	Source string
}

func (Comment) Node() string           { return "Comment" }
func (Comment) blockChild()            {}
func (Comment) exprListChild()         {}
func (Comment) fieldListChild()        {}
func (Comment) fileChild()             {}
func (Comment) identListChild()        {}
func (Comment) importListChild()       {}
func (Comment) paramListChild()        {}
func (Comment) typeListChild()         {}
func (x Comment) Pos() source.Position { return x.At }
func (x Comment) End() source.Position { return position.After(x.At, x.Source) }
func (x Comment) Dump() string         { return "Comment{" + x.Source + "}" }

func IsComment(node Node) bool {
	_, ok := node.(Comment)
	return ok
}

type ConstantDef struct {
	At        source.Position
	Public    bool
	ConstName string
	Value     ExprChild
	EndAt     source.Position
}

func (ConstantDef) Node() string           { return "ConstantDef" }
func (ConstantDef) fileChild()             {}
func (x ConstantDef) Pos() source.Position { return x.At }
func (x ConstantDef) End() source.Position { return x.EndAt }
func (x ConstantDef) Name() string         { return x.ConstName }
func (x ConstantDef) IsPublic() bool       { return x.Public }

func (x ConstantDef) Dump() string {
	s := x.ConstName + " = " + x.Value.Dump()
	if x.Public {
		s = "pub " + s
	}
	return "ConstantDef{" + s + "}"
}

type FunctionDef struct {
	At           source.Position
	Public       bool
	ReceiverName string
	ReceiverType *TypeSpec
	FuncName     string
	Params       []ParamListChild
	ParamsEndAt  source.Position
	Results      []TypeListChild
	BodyAt       source.Position
	Body         []BlockChild
	EndAt        source.Position
}

func (FunctionDef) Node() string           { return "FunctionDef" }
func (FunctionDef) fileChild()             {}
func (x FunctionDef) Pos() source.Position { return x.At }
func (x FunctionDef) End() source.Position { return x.EndAt }
func (x FunctionDef) Name() string         { return x.FuncName }
func (x FunctionDef) IsPublic() bool       { return x.Public }

func (x FunctionDef) Dump() string {
	s := x.FuncName + "("
	if x.Public {
		s = "pub " + s
	}

	delim := false
	for _, node := range x.Params {
		VisitParamListChild(node,
			func(Comment) {},
			func(node Parameter) {
				if delim {
					s += ", "
				}
				s += node.Dump()
				delim = true
			},
		)
	}

	switch len(x.Results) {
	case 0:
		s += ") "
	case 1:
		s += ") " + x.Results[0].Dump() + " "
	default:
		s += ") ("
		delim := false
		for _, node := range x.Results {
			VisitTypeListChild(node,
				func(Comment) {},
				func(node TypeSpec) {
					if delim {
						s += ", "
					}
					s += node.Dump()
					delim = true
				},
			)
		}
		s += ") "
	}

	return "FunctionDef{" + s + Block{x.At, x.Body, x.EndAt}.Dump() + "}"
}

type Field struct {
	At        source.Position
	FieldName string
	Type      TypeSpec
	Access    field.Access
	EndAt     source.Position
}

func (Field) Node() string           { return "Field" }
func (Field) fieldListChild()        {}
func (x Field) Pos() source.Position { return x.At }
func (x Field) End() source.Position { return x.EndAt }
func (x Field) Name() string         { return x.FieldName }
func (x Field) Dump() string         { return "Field{" + strings.Join(x.dumpRow(), " ") + "}" }

func (x Field) dumpRow() []string {
	if x.Access == field.AccessHidden {
		return []string{x.FieldName, x.Type.Dump()}
	}
	return []string{x.FieldName, x.Type.Dump(), x.Access.String()}
}

type Import struct {
	At    source.Position
	Path  string
	Names []IdentListChild
	EndAt source.Position
}

func (Import) Node() string           { return "Import" }
func (Import) blockChild()            {}
func (Import) fieldListChild()        {}
func (Import) fileChild()             {}
func (Import) importListChild()       {}
func (x Import) Pos() source.Position { return x.At }
func (x Import) End() source.Position { return x.EndAt }
func (x Import) Dump() string         { return "Import{" + x.stringInList() + "}" }

func (x Import) stringInList() string {
	s := x.Path

	delim := " ("
	for _, node := range x.Names {
		VisitIdentListChild(node,
			func(Comment) {},
			func(node Identifier) {
				s += delim + node.Name.String()
				delim = ", "
			},
		)
	}
	if s != x.Path {
		s += ")"
	}

	return s
}

type Imports struct {
	At      source.Position
	Imports []ImportListChild
	EndAt   source.Position
}

func (Imports) Node() string           { return "Imports" }
func (Imports) fileChild()             {}
func (x Imports) Pos() source.Position { return x.At }
func (x Imports) End() source.Position { return x.EndAt }

func (x Imports) Dump() string {
	s := "Imports{"

	delim := false
	for _, node := range x.Imports {
		VisitImportListChild(node,
			func(Comment) {},
			func(node Import) {
				if delim {
					s += "; "
				}
				s += node.stringInList()
				delim = true
			},
		)
	}

	return s + "}"
}

type Parameter struct {
	At        source.Position
	ParamName string
	Type      TypeSpec
	EndAt     source.Position
}

func (Parameter) Node() string           { return "Parameter" }
func (Parameter) paramListChild()        {}
func (x Parameter) Pos() source.Position { return x.At }
func (x Parameter) End() source.Position { return x.EndAt }
func (x Parameter) Dump() string         { return "Parameter{" + strings.Join(x.dumpRow(), " ") + "}" }
func (x Parameter) dumpRow() []string    { return []string{x.ParamName, x.Type.Dump()} }
func (x Parameter) Name() string         { return x.ParamName }

type TypeDef struct {
	At       source.Position
	Public   bool
	TypeName string
	Fields   []FieldListChild
	EndAt    source.Position
}

func (TypeDef) Node() string           { return "TypeDef" }
func (TypeDef) fileChild()             {}
func (x TypeDef) Pos() source.Position { return x.At }
func (x TypeDef) End() source.Position { return x.EndAt }
func (x TypeDef) Name() string         { return x.TypeName }
func (x TypeDef) IsPublic() bool       { return x.Public }

func (x TypeDef) Dump() string {
	s := x.TypeName + " {"
	if x.Public {
		s = "pub " + s
	}

	delim := false
	for _, node := range x.Fields {
		VisitFieldListChild(node,
			func(Comment) {},
			func(node Field) {
				if delim {
					s += "; "
				}
				s += node.Dump()
				delim = true
			},
			func(Import) {},
		)
	}

	return "TypeDef{" + s + "}}"
}
