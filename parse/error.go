// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
)

func newError(pos source.Position, msg string) error {
	return position.NewError(pos, msg)
}
