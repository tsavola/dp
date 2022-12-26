// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package position

import (
	"unicode/utf8"

	"github.com/tsavola/dp/source"
)

func After(p source.Position, s string) source.Position {
	for len(s) > 0 {
		c, n := utf8.DecodeRuneInString(s)
		s = s[n:]

		if c == '\n' {
			p.Line++
			p.Column = 1
		} else {
			p.Column++
		}
		p.ByteOffset += n
	}

	return p
}
