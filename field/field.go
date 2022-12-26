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
		return "hidden"
	case AccessVisible:
		return "visible"
	case AccessMutable:
		return "mutable"
	case AccessAssignable:
		return "assignable"
	default:
		return "<invalid access mode>"
	}
}
