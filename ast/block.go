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
func (x Expression) Dump() string            { return "Expression{" + x.Expr.Dump() + "}" }

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

func (x Assign) Dump() string {
	s := "Assign{"

	delim := ""
	for _, node := range x.Objects {
		s += delim + node.Dump()
		delim = ", "
	}

	delim = " = "
	for _, node := range x.Subjects {
		if !IsComment(node) {
			s += delim + node.Dump()
			delim = ", "
		}
	}

	return s + "}"
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

func (x Block) Dump() string {
	s := "Block{"

	delim := false
	for _, node := range x.Body {
		if !IsComment(node) {
			if delim {
				s += "; "
			}
			s += node.Dump()
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
func (Break) Dump() string              { return "Break" }

type Continue struct {
	source.Position
	End source.Position
}

func (Continue) Node() string              { return "continue" }
func (Continue) blockChild()               {}
func (x Continue) Pos() source.Position    { return x.Position }
func (x Continue) EndPos() source.Position { return x.End }
func (Continue) Dump() string              { return "Continue" }

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

func (x For) Dump() string {
	s := "For{"
	if x.Test != nil {
		s += x.Test.Dump() + " "
	}
	return s + Block{x.Position, x.Body, x.End}.Dump() + "}"
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

func (x If) Dump() string {
	s := "If{" + x.Test.Dump() + " " + Block{x.Position, x.Then, x.ThenEnd}.Dump()
	if len(x.Else) > 0 {
		s += " else " + Block{x.ThenEnd, x.Else, x.End}.Dump()
	}
	return s + "}"
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

func (x Return) Dump() string {
	s := "Return{"

	delim := " "
	for _, node := range x.Values {
		VisitExprListChild(node,
			func(node AssignerDereference) {
				s += delim + node.Dump()
				delim = ", "
			},
			func(Comment) {},
			func(node Expression) {
				s += delim + node.Dump()
				delim = ", "
			},
		)
	}

	return s + "}"
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
func (x VariableDecl) Dump() string            { return "VariableDecl{" + strings.Join(x.dumpRow(), " ") + "}" }

func (x VariableDecl) dumpRow() []string {
	typeName := "auto"
	if x.Type != nil {
		typeName = x.Type.Dump()
	}

	return []string{dumpNameList(x.Names), ":", typeName}
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
func (x VariableDef) Dump() string            { return "VariableDef{" + strings.Join(x.dumpRow(), " ") + "}" }

func (x VariableDef) dumpRow() []string {
	var exprs string
	var delim string

	for _, node := range x.Values {
		VisitExprListChild(node,
			func(node AssignerDereference) {
				exprs += delim + node.Dump()
				delim = ", "
			},
			func(Comment) {},
			func(node Expression) {
				exprs += delim + node.Dump()
				delim = ", "
			},
		)
	}

	return []string{dumpNameList(x.Names), ":=", exprs}
}

func dumpNameList(names []string) string {
	var s string
	var delim string

	for _, name := range names {
		s += delim + name
		delim = ", "
	}

	return s
}
