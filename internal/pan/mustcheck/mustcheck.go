// Copyright (c) 2025 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package mustcheck

import (
	"github.com/tsavola/dp/internal/pan"
)

var Check = pan.Check

func Must[T any](x T, err error) T {
	Check(err)
	return x
}
