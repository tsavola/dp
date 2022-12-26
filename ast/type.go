// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"slices"

	"github.com/tsavola/dp/source"
)

type Type struct {
	Assigner  bool
	Pointer   bool
	Reference bool
	Shared    bool
	Item      *Type         // Non-nil if array.
	Name      QualifiedName // Empty if array.
}

func (t Type) Equal(other Type) bool {
	if t.Item != nil {
		if other.Item == nil || !t.Item.Equal(*other.Item) {
			return false
		}
	} else {
		if other.Item != nil {
			return false
		}
	}
	return slices.Equal(t.Name, other.Name) && t.Shared == other.Shared && t.Reference == other.Reference && t.Pointer == other.Pointer && t.Assigner == other.Assigner
}

func (t Type) String() string {
	var s string
	if t.Item != nil {
		s = "[" + t.Item.String() + "]"
	} else {
		s = t.Name.String()
	}
	if t.Shared {
		s = "#" + s
	}
	if t.Reference {
		s = "&" + s
	}
	if t.Pointer {
		s = "*" + s
	}
	if t.Assigner {
		s = "=" + s
	}
	return s
}

type TypeSpec struct {
	source.Position
	Type
	End source.Position
}

func (TypeSpec) Node() string              { return "type specification" }
func (TypeSpec) typeListChild()            {}
func (x TypeSpec) Pos() source.Position    { return x.Position }
func (x TypeSpec) EndPos() source.Position { return x.End }
func (x TypeSpec) String() string          { return x.Type.String() }
