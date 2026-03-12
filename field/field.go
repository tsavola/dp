// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package field describes structure field access modes.
package field

type Access int

const (
	AccessHidden Access = iota
	AccessVisible
	AccessMutable
	AccessAssignable
)

func (a Access) String() string {
	switch a {
	case AccessHidden:
		return "Hidden"
	case AccessVisible:
		return "Visible"
	case AccessMutable:
		return "Mutable"
	case AccessAssignable:
		return "Assignable"
	default:
		return "<invalid access mode>"
	}
}
