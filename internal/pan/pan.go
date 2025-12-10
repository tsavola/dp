// Copyright (c) 2025 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package pan

import (
	"import.name/pan"
)

var z = new(pan.Zone)

var (
	Check   = z.Check
	Error   = z.Error
	Panic   = z.Panic
	Recover = z.Recover
	Wrap    = z.Wrap
)
