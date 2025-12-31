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
	source.Position
	Name QualifiedName
	End  source.Position
}

func (Identifier) Node() string              { return "identifier" }
func (Identifier) identListChild()           {}
func (x Identifier) Pos() source.Position    { return x.Position }
func (x Identifier) EndPos() source.Position { return x.End }
func (x Identifier) Dump() string            { return "Identifier{" + x.Name.String() + "}" }
