// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"strings"

	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
)

type Address struct {
	At    source.Position
	Expr  ExprChild
	EndAt source.Position
}

func (Address) Node() string           { return "Address" }
func (Address) exprChild()             {}
func (x Address) Pos() source.Position { return x.At }
func (x Address) End() source.Position { return x.EndAt }
func (x Address) Dump() string         { return "Address{" + x.Expr.Dump() + "}" }

type AssignerDereference struct {
	At    source.Position
	Name  string
	EndAt source.Position
}

func (AssignerDereference) Node() string           { return "AssignerDereference" }
func (AssignerDereference) assignListChild()       {}
func (AssignerDereference) exprListChild()         {}
func (x AssignerDereference) Pos() source.Position { return x.At }
func (x AssignerDereference) End() source.Position { return x.EndAt }
func (x AssignerDereference) Dump() string         { return "AssignerDereference{" + x.Name + "}" }

type Binary struct {
	Left  ExprChild
	Op    BinaryOp
	Right ExprChild
	EndAt source.Position
}

func (Binary) Node() string           { return "Binary" }
func (Binary) exprChild()             {}
func (x Binary) Pos() source.Position { return x.Left.Pos() }
func (x Binary) End() source.Position { return x.EndAt }

func (x Binary) Dump() string {
	return "Binary{" + x.Left.Dump() + " " + x.Op.Dump() + " " + x.Right.Dump() + "}"
}

type Boolean struct {
	At     source.Position
	Source string
}

func (Boolean) Node() string           { return "Boolean" }
func (Boolean) exprChild()             {}
func (x Boolean) Pos() source.Position { return x.At }
func (x Boolean) End() source.Position { return position.After(x.At, x.Source) }
func (x Boolean) Dump() string         { return "Boolean{" + x.Source + "}" }

type Call struct {
	Name  Selector
	Args  []ExprListChild
	EndAt source.Position
}

func (Call) Node() string           { return "Call" }
func (Call) assignListChild()       {}
func (Call) exprChild()             {}
func (x Call) Pos() source.Position { return x.Name.Pos() }
func (x Call) End() source.Position { return x.EndAt }

func (x Call) Dump() string {
	s := "Call{" + x.Name.Dump() + " ("

	delim := ""
	for _, node := range x.Args {
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
	return s + ")}"
}

type Character struct {
	At     source.Position
	Source string
}

func (Character) Node() string           { return "Character" }
func (Character) exprChild()             {}
func (x Character) Pos() source.Position { return x.At }
func (x Character) End() source.Position { return position.After(x.At, x.Source) }
func (x Character) Dump() string         { return "Character{" + x.Source + "}" }

type Clone struct {
	At    source.Position
	Expr  ExprChild
	EndAt source.Position
}

func (Clone) Node() string           { return "Clone" }
func (Clone) exprChild()             {}
func (x Clone) Pos() source.Position { return x.At }
func (x Clone) End() source.Position { return x.EndAt }
func (x Clone) Dump() string         { return "Clone{" + x.Expr.Dump() + "}" }

type Index struct {
	Name  Selector
	Index ExprChild
	EndAt source.Position
}

func (Index) Node() string           { return "Index" }
func (Index) assignListChild()       {}
func (Index) exprChild()             {}
func (x Index) Pos() source.Position { return x.Name.Pos() }
func (x Index) End() source.Position { return x.EndAt }
func (x Index) Dump() string         { return "Index{" + x.Name.Dump() + " [" + x.Index.Dump() + "]}" }

type Integer struct {
	At     source.Position
	Source string
}

func (Integer) Node() string           { return "Integer" }
func (Integer) exprChild()             {}
func (x Integer) Pos() source.Position { return x.At }
func (x Integer) End() source.Position { return position.After(x.At, x.Source) }
func (x Integer) Dump() string         { return "Integer{" + x.Source + "}" }

type Nil struct {
	At    source.Position
	EndAt source.Position
}

func (Nil) Node() string           { return "Nil" }
func (Nil) exprChild()             {}
func (x Nil) Pos() source.Position { return x.At }
func (x Nil) End() source.Position { return x.EndAt }
func (Nil) Dump() string           { return "Nil" }

type PointerDereference struct {
	At    source.Position
	Expr  ExprChild
	EndAt source.Position
}

func (PointerDereference) Node() string           { return "PointerDereference" }
func (PointerDereference) exprChild()             {}
func (x PointerDereference) Pos() source.Position { return x.At }
func (x PointerDereference) End() source.Position { return x.EndAt }
func (x PointerDereference) Dump() string         { return "PointerDereference{" + x.Expr.Dump() + "}" }

type Selector struct {
	At    source.Position
	Name  []string
	EndAt source.Position
}

func (Selector) Node() string           { return "Selector" }
func (Selector) assignListChild()       {}
func (Selector) exprChild()             {}
func (x Selector) Pos() source.Position { return x.At }
func (x Selector) End() source.Position { return x.EndAt }
func (x Selector) Dump() string         { return "Selector{" + x.String() + "}" }
func (x Selector) String() string       { return strings.Join(x.Name, ".") }

type String struct {
	At     source.Position
	Source string
}

func (String) Node() string           { return "String" }
func (String) exprChild()             {}
func (x String) Pos() source.Position { return x.At }
func (x String) End() source.Position { return position.After(x.At, x.Source) }
func (x String) Dump() string         { return "String{" + x.Source + "}" }

type Unary struct {
	At    source.Position
	Op    UnaryOp
	Expr  ExprChild
	EndAt source.Position
}

func (Unary) Node() string           { return "Unary" }
func (Unary) exprChild()             {}
func (x Unary) Pos() source.Position { return x.At }
func (x Unary) End() source.Position { return x.EndAt }
func (x Unary) Dump() string         { return "Unary{" + x.Op.Dump() + " " + x.Expr.Dump() + "}" }

type Zero struct {
	At    source.Position
	EndAt source.Position
}

func (Zero) Node() string           { return "Zero" }
func (Zero) exprChild()             {}
func (x Zero) Pos() source.Position { return x.At }
func (x Zero) End() source.Position { return x.EndAt }
func (Zero) Dump() string           { return "Zero" }
