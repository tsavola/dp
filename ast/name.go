// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package ast

import (
	"strings"

	"github.com/tsavola/dp/source"
)

type QualifiedName []string

func (name QualifiedName) Namespace() string { return strings.Join(name[:len(name)-1], "::") }
func (name QualifiedName) Short() string     { return name[len(name)-1] }
func (name QualifiedName) String() string    { return strings.Join(name, "::") }

type Identifier struct {
	At    source.Position
	Name  QualifiedName
	EndAt source.Position
}

func (Identifier) Node() string           { return "Identifier" }
func (Identifier) identListChild()        {}
func (x Identifier) Pos() source.Position { return x.At }
func (x Identifier) End() source.Position { return x.EndAt }
func (x Identifier) Dump() string         { return "Identifier{" + x.Name.String() + "}" }
