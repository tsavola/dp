// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package lex

import (
	"unicode/utf8"

	"github.com/tsavola/dp/internal/pan"
	"github.com/tsavola/dp/source"
)

// scan state.
type scan struct {
	source.Position
	text string
}

func (s scan) pos() source.Position {
	return s.Position
}

// peek returns 0 on EOF.  Zero code point is returned as utf8.RuneError.
func (s *scan) peek() rune {
	r, _ := s.peekSize()
	return r
}

// advance must only be called after peek returned a nonzero rune.
func (s *scan) advance() {
	c, n := s.peekSize()
	if c == 0 {
		panic("advanced on EOF")
	}

	if c == '\n' {
		s.Line++
		s.Column = 1
	} else {
		s.Column++
	}
	s.ByteOffset += n
}

// peekSize returns 0 on EOF.  Zero code point is returned as utf8.RuneError.
func (s *scan) peekSize() (rune, int) {
	if s.ByteOffset == len(s.text) {
		return 0, 0
	}

	c, n := utf8.DecodeRuneInString(s.text[s.ByteOffset:])
	if c == utf8.RuneError {
		pan.Panic(decodeError(s.pos()))
	}
	if c == 0 {
		return utf8.RuneError, n
	}

	return c, n
}

// until returns a part of the source.
func (s *scan) until(endByteOffset int) string {
	return s.text[s.ByteOffset:endByteOffset]
}
