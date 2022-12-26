// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package lex

import (
	"github.com/tsavola/dp/internal/position"
)

func tokenError(s scan) error {
	return position.NewError(s.pos(), "illegal token")
}
