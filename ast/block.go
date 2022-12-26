// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"strings"

	"github.com/tsavola/dp/source"
)

type Expression struct {
	Expr ExprChild
}

func (Expression) Node() string              { return "expression" }
func (Expression) blockChild()               {}
func (Expression) exprListChild()            {}
func (x Expression) Pos() source.Position    { return x.Expr.Pos() }
func (x Expression) EndPos() source.Position { return x.Expr.EndPos() }
func (x Expression) String() string          { return x.Expr.String() }

type Assign struct {
	source.Position
	Objects  []AssignListChild
	Subjects []ExprListChild
	End      source.Position
}

func (Assign) Node() string              { return "assignment" }
func (Assign) blockChild()               {}
func (x Assign) Pos() source.Position    { return x.Position }
func (x Assign) EndPos() source.Position { return x.End }

func (x Assign) String() string {
	var s string

	delim := ""
	for _, node := range x.Objects {
		s += delim + node.String()
		delim = ", "
	}

	delim = " = "
	for _, node := range x.Subjects {
		if !IsComment(node) {
			s += delim + node.String()
			delim = ", "
		}
	}

	return s
}

type Block struct {
	source.Position
	Body []BlockChild
	End  source.Position
}

func (Block) Node() string              { return "block" }
func (Block) blockChild()               {}
func (x Block) Pos() source.Position    { return x.Position }
func (x Block) EndPos() source.Position { return x.End }

func (x Block) String() string {
	s := "{"

	delim := false
	for _, node := range x.Body {
		if !IsComment(node) {
			if delim {
				s += "; "
			}
			s += node.String()
			delim = true
		}
	}

	return s + "}"
}

type Break struct {
	source.Position
	End source.Position
}

func (Break) Node() string              { return "break" }
func (Break) blockChild()               {}
func (x Break) Pos() source.Position    { return x.Position }
func (x Break) EndPos() source.Position { return x.End }
func (Break) String() string            { return "break" }

type Continue struct {
	source.Position
	End source.Position
}

func (Continue) Node() string              { return "continue" }
func (Continue) blockChild()               {}
func (x Continue) Pos() source.Position    { return x.Position }
func (x Continue) EndPos() source.Position { return x.End }
func (Continue) String() string            { return "continue" }

type For struct {
	source.Position
	Test    ExprChild // Nil if infinite loop.
	BodyPos source.Position
	Body    []BlockChild
	End     source.Position
}

func (For) Node() string              { return "for statement" }
func (For) blockChild()               {}
func (x For) Pos() source.Position    { return x.Position }
func (x For) EndPos() source.Position { return x.End }

func (x For) String() string {
	s := "for "
	if x.Test != nil {
		s += x.Test.String() + " "
	}
	return s + Block{x.Position, x.Body, x.End}.String()
}

type If struct {
	source.Position
	Test    ExprChild
	ThenPos source.Position
	Then    []BlockChild
	ThenEnd source.Position
	Else    []BlockChild
	End     source.Position
}

func (If) Node() string              { return "if statement" }
func (If) blockChild()               {}
func (x If) Pos() source.Position    { return x.Position }
func (x If) EndPos() source.Position { return x.End }

func (x If) String() string {
	s := "if " + x.Test.String() + " " + Block{x.Position, x.Then, x.ThenEnd}.String()
	if len(x.Else) > 0 {
		s += " else " + Block{x.ThenEnd, x.Else, x.End}.String()
	}
	return s
}

type Return struct {
	source.Position
	Values []ExprListChild
	End    source.Position
}

func (Return) Node() string              { return "return statement" }
func (Return) blockChild()               {}
func (x Return) Pos() source.Position    { return x.Position }
func (x Return) EndPos() source.Position { return x.End }

func (x Return) String() string {
	s := "return"

	delim := " "
	for _, node := range x.Values {
		VisitExprListChild(node,
			func(node AssignerDereference) {
				s += delim + node.String()
				delim = ", "
			},
			func(Comment) {},
			func(node Expression) {
				s += delim + node.String()
				delim = ", "
			},
		)
	}

	return s
}

type VariableDecl struct {
	source.Position
	Names []string
	Type  *TypeSpec // Nil if auto.
	End   source.Position
}

func (VariableDecl) Node() string              { return "variable declaration" }
func (VariableDecl) blockChild()               {}
func (x VariableDecl) Pos() source.Position    { return x.Position }
func (x VariableDecl) EndPos() source.Position { return x.End }
func (x VariableDecl) String() string          { return strings.Join(x.row(), " ") }

func (x VariableDecl) row() []string {
	typeName := "auto"
	if x.Type != nil {
		typeName = x.Type.String()
	}

	return []string{formatNameList(x.Names), ":", typeName}
}

type VariableDef struct {
	source.Position
	Names  []string
	Values []ExprListChild
	End    source.Position
}

func (VariableDef) Node() string              { return "variable definition" }
func (VariableDef) blockChild()               {}
func (x VariableDef) Pos() source.Position    { return x.Position }
func (x VariableDef) EndPos() source.Position { return x.End }
func (x VariableDef) String() string          { return strings.Join(x.row(), " ") }

func (x VariableDef) row() []string {
	var exprs string
	var delim string

	for _, node := range x.Values {
		VisitExprListChild(node,
			func(node AssignerDereference) {
				exprs += delim + node.String()
				delim = ", "
			},
			func(Comment) {},
			func(node Expression) {
				exprs += delim + node.String()
				delim = ", "
			},
		)
	}

	return []string{formatNameList(x.Names), ":=", exprs}
}

func formatNameList(names []string) string {
	var s string
	var delim string

	for _, name := range names {
		s += delim + name
		delim = ", "
	}

	return s
}
