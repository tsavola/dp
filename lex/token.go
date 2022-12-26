// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package lex implements lexical analysis.
package lex

import (
	"unicode"

	"github.com/tsavola/dp/source"
	"github.com/tsavola/dp/token"
	"import.name/pan"
)

func File(pos source.Position, text string) ([]token.Token, error) {
	var tokens []token.Token

	err := pan.Recover(func() {
		tokens = tokenize(scan{pos, text})
	})

	return tokens, err
}

func tokenize(s scan) (tokens []token.Token) {
next:
	for s.peek() != 0 {
		for _, f := range []func(scan) (scan, token.Token){
			tokenizeSpace,
			tokenizer(token.Newline, "\n", nil),
			tokenizeComment,

			tokenizer(token.Auto, "auto", wordRune),
			tokenizer(token.Break, "break", wordRune),
			tokenizer(token.Clone, "clone", wordRune),
			tokenizer(token.Continue, "continue", wordRune),
			tokenizer(token.Else, "else", wordRune),
			tokenizer(token.False, "false", wordRune),
			tokenizer(token.For, "for", wordRune),
			tokenizer(token.If, "if", wordRune),
			tokenizer(token.Import, "import", wordRune),
			tokenizer(token.Nil, "nil", wordRune),
			tokenizer(token.Return, "return", wordRune),
			tokenizer(token.True, "true", wordRune),

			tokenizeWord,

			tokenizeInteger,
			quoteTokenizer(token.Character, '\''),
			quoteTokenizer(token.String, '"'),
			quoteTokenizer(token.String, '`'),

			tokenizer(token.Plus, "+", nil),
			tokenizer(token.Minus, "-", nil),
			tokenizer(token.Asterisk, "*", nil),
			tokenizer(token.Slash, "/", nil),
			tokenizer(token.Percent, "%", nil),

			tokenizer(token.LogicalAnd, "&&", nil),
			tokenizer(token.LogicalOr, "||", nil),

			tokenizer(token.AndNot, "&^", nil),
			tokenizer(token.Ampersand, "&", nil),
			tokenizer(token.Pipe, "|", nil),
			tokenizer(token.Caret, "^", nil),

			tokenizer(token.ShiftLeft, "<<", nil),
			tokenizer(token.ShiftRight, ">>", nil),

			tokenizer(token.Equal, "==", nil),
			tokenizer(token.NotEqual, "!=", nil),
			tokenizer(token.LessOrEqual, "<=", nil),
			tokenizer(token.GreaterOrEqual, ">=", nil),
			tokenizer(token.Less, "<", nil),
			tokenizer(token.Greater, ">", nil),

			tokenizer(token.Exclamation, "!", nil),

			tokenizer(token.Assign, "=", nil),
			tokenizer(token.Define, ":=", nil),

			tokenizer(token.Comma, ",", nil),
			tokenizer(token.Period, ".", nil),
			tokenizer(token.Semicolon, ";", nil),
			tokenizer(token.Colons, "::", nil),
			tokenizer(token.Colon, ":", nil),
			tokenizer(token.Hash, "#", nil),

			tokenizer(token.ParenLeft, "(", nil),
			tokenizer(token.BracketLeft, "[", nil),
			tokenizer(token.BraceLeft, "{", nil),

			tokenizer(token.ParenRight, ")", nil),
			tokenizer(token.BracketRight, "]", nil),
			tokenizer(token.BraceRight, "}", nil),
		} {
			if pan.Recover(func() {
				after, t := f(s)
				tokens = append(tokens, t)
				s = after
			}) == nil {
				continue next
			}
		}

		pan.Panic(tokenError(s))
	}

	return
}

func tokenizeSpace(s scan) (scan, token.Token) {
	start := s

	for {
		if c := s.peek(); c == '\n' || !unicode.IsSpace(c) {
			return makeToken(token.Space, start, s)
		}
		s.advance()
	}
}

func tokenizeComment(s scan) (scan, token.Token) {
	start := s

	for range 2 {
		if s.peek() != '/' {
			pan.Panic(tokenError(start))
		}
		s.advance()
	}

	for {
		switch s.peek() {
		case '\n', 0:
			return makeToken(token.Comment, start, s)
		}
		s.advance()
	}
}

func tokenizeWord(s scan) (scan, token.Token) {
	start := s

	if !wordStartRune(s.peek()) {
		pan.Panic(tokenError(start))
	}

	for {
		s.advance()

		if !wordRune(s.peek()) {
			return makeToken(token.Word, start, s)
		}
	}
}

func tokenizeInteger(s scan) (scan, token.Token) {
	start := s

	for {
		if !unicode.IsDigit(s.peek()) {
			return makeToken(token.Integer, start, s)
		}
		s.advance()
	}
}

func quoteTokenizer(k token.Kind, quote rune) func(scan) (scan, token.Token) {
	return func(s scan) (scan, token.Token) {
		start := s

		if s.peek() != quote {
			pan.Panic(tokenError(start))
		}
		s.advance()

		for {
			c := s.peek()
			if c == 0 {
				pan.Panic(tokenError(start))
			}

			s.advance()

			switch c {
			case '\\':
				if s.peek() == 0 {
					pan.Panic(tokenError(start))
				}
				s.advance()

			case quote:
				return makeToken(k, start, s)
			}
		}
	}
}

func tokenizer(k token.Kind, source string, failIfTrailing func(rune) bool) func(scan) (scan, token.Token) {
	return func(s scan) (scan, token.Token) {
		start := s

		for _, wanted := range source {
			switch s.peek() {
			case 0:
				pan.Panic(tokenError(start))

			case wanted:
				s.advance()

			default:
				pan.Panic(tokenError(start))
			}
		}

		if failIfTrailing != nil && failIfTrailing(s.peek()) {
			pan.Panic(tokenError(start))
		}

		return makeToken(k, start, s)
	}
}

func makeToken(k token.Kind, start, end scan) (scan, token.Token) {
	source := start.until(end.ByteOffset)
	if source == "" {
		pan.Panic(tokenError(start))
	}

	return end, token.Token{start.pos(), k, source}
}

func wordStartRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func wordRune(r rune) bool {
	return wordStartRune(r) || unicode.IsDigit(r)
}
