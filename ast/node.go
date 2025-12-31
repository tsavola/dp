// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package ast contains abstract syntax tree node types.
package ast

import (
	"github.com/tsavola/dp/source"
)

type Node interface {
	Node() string
	Pos() source.Position
	EndPos() source.Position
	Dump() string
}

type AssignListChild interface {
	Node
	assignListChild()
}

type BlockChild interface {
	Node
	blockChild()
}

type ExprChild interface {
	Node
	exprChild()
}

type ExprListChild interface {
	Node
	exprListChild()
}

type FieldListChild interface {
	Node
	fieldListChild()
}

type FileChild interface {
	Node
	fileChild()
}

type IdentListChild interface {
	Node
	identListChild()
}

type ImportListChild interface {
	Node
	importListChild()
}

type ParamListChild interface {
	Node
	paramListChild()
}

type TypeListChild interface {
	Node
	typeListChild()
}
