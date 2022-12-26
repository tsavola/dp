// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"strings"

	"github.com/tsavola/dp/source"
)

type Address struct {
	source.Position
	Expr ExprChild
	End  source.Position
}

func (Address) Node() string              { return "address operation" }
func (Address) exprChild()                {}
func (x Address) Pos() source.Position    { return x.Position }
func (x Address) EndPos() source.Position { return x.End }
func (x Address) String() string          { return "&" + x.Expr.String() }

type AssignerDereference struct {
	source.Position
	Name string
	End  source.Position
}

func (AssignerDereference) Node() string              { return "assigner dereference" }
func (AssignerDereference) assignListChild()          {}
func (AssignerDereference) exprListChild()            {}
func (x AssignerDereference) Pos() source.Position    { return x.Position }
func (x AssignerDereference) EndPos() source.Position { return x.End }
func (x AssignerDereference) String() string          { return "(" + x.Name + ")" }

type Binary struct {
	Left  ExprChild
	Op    BinaryOp
	Right ExprChild
	End   source.Position
}

func (Binary) Node() string              { return "binary operation" }
func (Binary) exprChild()                {}
func (x Binary) Pos() source.Position    { return x.Left.Pos() }
func (x Binary) EndPos() source.Position { return x.End }

func (x Binary) String() string {
	return x.Left.String() + " " + x.Op.String() + " " + x.Right.String()
}

type Boolean struct {
	source.Position
	Value bool
	End   source.Position
}

func (Boolean) Node() string              { return "boolean literal" }
func (Boolean) exprChild()                {}
func (x Boolean) Pos() source.Position    { return x.Position }
func (x Boolean) EndPos() source.Position { return x.End }

func (x Boolean) String() string {
	if x.Value {
		return "true"
	}
	return "false"
}

type Call struct {
	Name Selector
	Args []ExprListChild
	End  source.Position
}

func (Call) Node() string              { return "function call" }
func (Call) assignListChild()          {}
func (Call) exprChild()                {}
func (x Call) Pos() source.Position    { return x.Name.Pos() }
func (x Call) EndPos() source.Position { return x.End }

func (x Call) String() string {
	s := x.Name.String() + "("

	delim := ""
	for _, node := range x.Args {
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
	return s + ")"
}

type Character struct {
	source.Position
	Source string
	End    source.Position
}

func (Character) Node() string              { return "character literal" }
func (Character) exprChild()                {}
func (x Character) Pos() source.Position    { return x.Position }
func (x Character) EndPos() source.Position { return x.End }
func (x Character) String() string          { return x.Source }

type Clone struct {
	source.Position
	Expr ExprChild
	End  source.Position
}

func (Clone) Node() string              { return "clone operation" }
func (Clone) exprChild()                {}
func (x Clone) Pos() source.Position    { return x.Position }
func (x Clone) EndPos() source.Position { return x.End }
func (x Clone) String() string          { return "clone " + x.Expr.String() }

type Index struct {
	Name  Selector
	Index ExprChild
	End   source.Position
}

func (Index) Node() string              { return "index operation" }
func (Index) assignListChild()          {}
func (Index) exprChild()                {}
func (x Index) Pos() source.Position    { return x.Name.Pos() }
func (x Index) EndPos() source.Position { return x.End }
func (x Index) String() string          { return x.Name.String() + "[" + x.Index.String() + "]" }

type Integer struct {
	source.Position
	Source string
	End    source.Position
}

func (Integer) Node() string              { return "integer literal" }
func (Integer) exprChild()                {}
func (x Integer) Pos() source.Position    { return x.Position }
func (x Integer) EndPos() source.Position { return x.End }
func (x Integer) String() string          { return x.Source }

type Nil struct {
	source.Position
	End source.Position
}

func (Nil) Node() string              { return "nil literal" }
func (Nil) exprChild()                {}
func (x Nil) Pos() source.Position    { return x.Position }
func (x Nil) EndPos() source.Position { return x.End }
func (Nil) String() string            { return "nil" }

type PointerDereference struct {
	source.Position
	Expr ExprChild
	End  source.Position
}

func (PointerDereference) Node() string              { return "pointer dereference" }
func (PointerDereference) exprChild()                {}
func (x PointerDereference) Pos() source.Position    { return x.Position }
func (x PointerDereference) EndPos() source.Position { return x.End }
func (x PointerDereference) String() string          { return "*" + x.Expr.String() }

type Selector struct {
	source.Position
	Name []string
	End  source.Position
}

func (Selector) Node() string              { return "selector operation" }
func (Selector) assignListChild()          {}
func (Selector) exprChild()                {}
func (x Selector) Pos() source.Position    { return x.Position }
func (x Selector) EndPos() source.Position { return x.End }
func (x Selector) String() string          { return strings.Join(x.Name, ".") }

type String struct {
	source.Position
	Source string
	End    source.Position
}

func (String) Node() string              { return "string literal" }
func (String) exprChild()                {}
func (x String) Pos() source.Position    { return x.Position }
func (x String) EndPos() source.Position { return x.End }
func (x String) String() string          { return x.Source }

type Unary struct {
	source.Position
	Op   UnaryOp
	Expr ExprChild
	End  source.Position
}

func (Unary) Node() string              { return "unary operation" }
func (Unary) exprChild()                {}
func (x Unary) Pos() source.Position    { return x.Position }
func (x Unary) EndPos() source.Position { return x.End }
func (x Unary) String() string          { return x.Op.String() + x.Expr.String() }

type Zero struct {
	source.Position
	End source.Position
}

func (Zero) Node() string              { return "zero literal" }
func (Zero) exprChild()                {}
func (x Zero) Pos() source.Position    { return x.Position }
func (x Zero) EndPos() source.Position { return x.End }
func (Zero) String() string            { return "{}" }
