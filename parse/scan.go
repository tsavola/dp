// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"github.com/tsavola/dp/internal/pan"
	"github.com/tsavola/dp/source"
	"github.com/tsavola/dp/token"
)

// scan state.
type scan struct {
	tokens []token.Token
	last   source.Position
}

func (s scan) pos() source.Position {
	s.peek() // Skip space.
	if len(s.tokens) == 0 {
		return source.Position{}
	}
	return s.tokens[0].Pos()
}

// peek at the next token.  Space tokens are skipped.  Last position is not
// updated.
func (s *scan) peek() token.Token {
	var zero token.Token

	for len(s.tokens) > 0 {
		if t := s.tokens[0]; t.Kind != token.Space {
			return t
		}
		s.tokens = s.tokens[1:]
	}

	return zero
}

// skip token if it's next.  Space tokens are skipped.  Last position is
// updated on success.
func (s *scan) skip(wanted token.Kind) bool {
	_, ok := s.skim(wanted)
	return ok
}

// skim returns token if it's next.  Space tokens are skipped.  Last position
// is updated on success.
func (s *scan) skim(wanted token.Kind) (token.Token, bool) {
	if s.peek().Kind != wanted {
		return token.Token{}, false
	}

	t := s.tokens[0]
	s.tokens = s.tokens[1:]

	if len(s.tokens) > 0 {
		s.last = s.pos()
	} else {
		s.last = t.EndPos()
	}

	return t, true
}

// take returns token or panics.  Space tokens are skipped.  Last position is
// updated on success.
func (s *scan) take(wanted token.Kind, errorMessage string) token.Token {
	t, ok := s.skim(wanted)
	if !ok {
		pan.Panic(newError(s.pos(), errorMessage))
	}
	return t
}

func peekEOF(s scan) (scan, bool) {
	s.peek() // Skip space.
	return s, len(s.tokens) == 0
}

func skipper(t token.Kind) func(scan) (scan, bool) {
	return func(s scan) (scan, bool) {
		ok := s.skip(t)
		return s, ok
	}
}
