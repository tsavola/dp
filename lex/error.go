// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package lex

import (
	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
)

func decodeError(pos source.Position) error {
	return position.NewError(pos, "invalid UTF-8 encoding")
}

func tokenError(s scan) error {
	return position.NewError(s.pos(), "illegal token")
}
