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

func (Expression) Node() string           { return "Expression" }
func (Expression) blockChild()            {}
func (Expression) exprListChild()         {}
func (x Expression) Pos() source.Position { return x.Expr.Pos() }
func (x Expression) End() source.Position { return x.Expr.End() }
func (x Expression) Dump() string         { return "Expression{" + x.Expr.Dump() + "}" }

type Assign struct {
	At      source.Position
	Targets []AssignListChild
	Values  []ExprListChild
	EndAt   source.Position
}

func (Assign) Node() string           { return "Assignment" }
func (Assign) blockChild()            {}
func (x Assign) Pos() source.Position { return x.At }
func (x Assign) End() source.Position { return x.EndAt }

func (x Assign) Dump() string {
	s := "Assign{"

	delim := ""
	for _, node := range x.Targets {
		s += delim + node.Dump()
		delim = ", "
	}

	delim = " = "
	for _, node := range x.Values {
		if !IsComment(node) {
			s += delim + node.Dump()
			delim = ", "
		}
	}

	return s + "}"
}

type Block struct {
	At    source.Position
	Body  []BlockChild
	EndAt source.Position
}

func (Block) Node() string           { return "Block" }
func (Block) blockChild()            {}
func (x Block) Pos() source.Position { return x.At }
func (x Block) End() source.Position { return x.EndAt }

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
	At    source.Position
	EndAt source.Position
}

func (Break) Node() string           { return "Break" }
func (Break) blockChild()            {}
func (x Break) Pos() source.Position { return x.At }
func (x Break) End() source.Position { return x.EndAt }
func (Break) Dump() string           { return "Break" }

type Continue struct {
	At    source.Position
	EndAt source.Position
}

func (Continue) Node() string           { return "Continue" }
func (Continue) blockChild()            {}
func (x Continue) Pos() source.Position { return x.At }
func (x Continue) End() source.Position { return x.EndAt }
func (Continue) Dump() string           { return "Continue" }

type For struct {
	At     source.Position
	Test   ExprChild // Nil if infinite loop.
	BodyAt source.Position
	Body   []BlockChild
	EndAt  source.Position
}

func (For) Node() string           { return "For" }
func (For) blockChild()            {}
func (x For) Pos() source.Position { return x.At }
func (x For) End() source.Position { return x.EndAt }

func (x For) Dump() string {
	s := "For{"
	if x.Test != nil {
		s += x.Test.Dump() + " "
	}
	return s + Block{x.At, x.Body, x.EndAt}.Dump() + "}"
}

type If struct {
	At        source.Position
	Test      ExprChild
	ThenAt    source.Position
	Then      []BlockChild
	ThenEndAt source.Position
	Else      []BlockChild
	EndAt     source.Position
}

func (If) Node() string           { return "If" }
func (If) blockChild()            {}
func (x If) Pos() source.Position { return x.At }
func (x If) End() source.Position { return x.EndAt }

func (x If) Dump() string {
	s := "If{" + x.Test.Dump() + " " + Block{x.At, x.Then, x.ThenEndAt}.Dump()
	if len(x.Else) > 0 {
		s += " else " + Block{x.ThenEndAt, x.Else, x.EndAt}.Dump()
	}
	return s + "}"
}

type Return struct {
	At     source.Position
	Values []ExprListChild
	EndAt  source.Position
}

func (Return) Node() string           { return "Return" }
func (Return) blockChild()            {}
func (x Return) Pos() source.Position { return x.At }
func (x Return) End() source.Position { return x.EndAt }

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
	At    source.Position
	Names []string
	Type  *TypeSpec // Nil if auto.
	EndAt source.Position
}

func (VariableDecl) Node() string           { return "VariableDecl" }
func (VariableDecl) blockChild()            {}
func (x VariableDecl) Pos() source.Position { return x.At }
func (x VariableDecl) End() source.Position { return x.EndAt }
func (x VariableDecl) Dump() string         { return "VariableDecl{" + strings.Join(x.dumpRow(), " ") + "}" }

func (x VariableDecl) dumpRow() []string {
	typeName := "auto"
	if x.Type != nil {
		typeName = x.Type.Dump()
	}

	return []string{dumpNameList(x.Names), ":", typeName}
}

type VariableDef struct {
	At     source.Position
	Names  []string
	Values []ExprListChild
	EndAt  source.Position
}

func (VariableDef) Node() string           { return "VariableDef" }
func (VariableDef) blockChild()            {}
func (x VariableDef) Pos() source.Position { return x.At }
func (x VariableDef) End() source.Position { return x.EndAt }
func (x VariableDef) Dump() string         { return "VariableDef{" + strings.Join(x.dumpRow(), " ") + "}" }

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
