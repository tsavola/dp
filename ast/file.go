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
	source.Position
	Source string
}

func (Comment) Node() string              { return "comment" }
func (Comment) blockChild()               {}
func (Comment) exprListChild()            {}
func (Comment) fieldListChild()           {}
func (Comment) fileChild()                {}
func (Comment) identListChild()           {}
func (Comment) importListChild()          {}
func (Comment) paramListChild()           {}
func (Comment) typeListChild()            {}
func (x Comment) Pos() source.Position    { return x.Position }
func (x Comment) EndPos() source.Position { return position.After(x.Position, x.Source) }
func (x Comment) Dump() string            { return "Comment{" + x.Source + "}" }

func IsComment(node Node) bool {
	_, ok := node.(Comment)
	return ok
}

type ConstantDef struct {
	source.Position
	Public    bool
	ConstName string
	Value     ExprChild
	End       source.Position
}

func (ConstantDef) Node() string              { return "constant definition" }
func (ConstantDef) fileChild()                {}
func (x ConstantDef) Pos() source.Position    { return x.Position }
func (x ConstantDef) EndPos() source.Position { return x.End }
func (x ConstantDef) Name() string            { return x.ConstName }

func (x ConstantDef) Dump() string {
	s := x.ConstName + " = " + x.Value.Dump()
	if x.Public {
		s = "pub " + s
	}
	return "ConstantDef{" + s + "}"
}

type FunctionDef struct {
	source.Position
	Public       bool
	ReceiverName string
	ReceiverType *TypeSpec
	FuncName     string
	Params       []ParamListChild
	ParamsEnd    source.Position
	Results      []TypeListChild
	BodyPos      source.Position
	Body         []BlockChild
	End          source.Position
}

func (FunctionDef) Node() string              { return "function definition" }
func (FunctionDef) fileChild()                {}
func (x FunctionDef) Pos() source.Position    { return x.Position }
func (x FunctionDef) EndPos() source.Position { return x.End }
func (x FunctionDef) Name() string            { return x.FuncName }

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

	return "FunctionDef{" + s + Block{x.Position, x.Body, x.End}.Dump() + "}"
}

type Field struct {
	source.Position
	FieldName string
	Type      TypeSpec
	Access    field.Access
	End       source.Position
}

func (Field) Node() string              { return "field" }
func (Field) fieldListChild()           {}
func (x Field) Pos() source.Position    { return x.Position }
func (x Field) EndPos() source.Position { return x.End }
func (x Field) Name() string            { return x.FieldName }
func (x Field) Dump() string            { return "Field{" + strings.Join(x.dumpRow(), " ") + "}" }

func (x Field) dumpRow() []string {
	if x.Access == field.AccessHidden {
		return []string{x.FieldName, x.Type.Dump()}
	}
	return []string{x.FieldName, x.Type.Dump(), x.Access.String()}
}

type Import struct {
	source.Position
	Path  string
	Names []IdentListChild
	End   source.Position
}

func (Import) Node() string              { return "import" }
func (Import) blockChild()               {}
func (Import) fieldListChild()           {}
func (Import) fileChild()                {}
func (Import) importListChild()          {}
func (x Import) Pos() source.Position    { return x.Position }
func (x Import) EndPos() source.Position { return x.End }
func (x Import) Dump() string            { return "Import{" + x.stringInList() + "}" }

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
	source.Position
	Imports []ImportListChild
	End     source.Position
}

func (Imports) Node() string              { return "import list" }
func (Imports) fileChild()                {}
func (x Imports) Pos() source.Position    { return x.Position }
func (x Imports) EndPos() source.Position { return x.End }

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
	source.Position
	ParamName string
	Type      TypeSpec
	End       source.Position
}

func (Parameter) Node() string              { return "parameter" }
func (Parameter) paramListChild()           {}
func (x Parameter) Pos() source.Position    { return x.Position }
func (x Parameter) EndPos() source.Position { return x.End }
func (x Parameter) Dump() string            { return "Parameter{" + strings.Join(x.dumpRow(), " ") + "}" }
func (x Parameter) dumpRow() []string       { return []string{x.ParamName, x.Type.Dump()} }
func (x Parameter) Name() string            { return x.ParamName }

type TypeDef struct {
	source.Position
	Public   bool
	TypeName string
	Fields   []FieldListChild
	End      source.Position
}

func (TypeDef) Node() string              { return "type definition" }
func (TypeDef) fileChild()                {}
func (x TypeDef) Pos() source.Position    { return x.Position }
func (x TypeDef) EndPos() source.Position { return x.End }
func (x TypeDef) Name() string            { return x.TypeName }

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
