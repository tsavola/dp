// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package token enumerates lexical tokens.
package token

import (
	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
)

type Token struct {
	source.Position
	Kind
	Source string
}

func (t Token) Pos() source.Position {
	return t.Position
}

// EndPos returns the position immediately after the token.
func (t Token) EndPos() source.Position {
	return position.After(t.Position, t.Source)
}

func (t Token) String() string {
	return t.Kind.String()
}

// Kind of token.
type Kind int

const (
	_ Kind = iota

	Space // Excluding newline.
	Newline
	Comment

	// Keywords
	Auto
	Break
	Clone
	Continue
	Else
	False
	For
	If
	Import
	Nil
	Return
	True

	// Identifier, or keyword at file level or in type definition
	Word

	// Basic type literals
	Integer   // 12345
	Character // 'a'
	String    // "abc" or `abc`

	// Operators and delimiters
	Plus     // +
	Minus    // -
	Asterisk // *
	Slash    // /
	Percent  // %

	LogicalAnd // &&
	LogicalOr  // ||

	AndNot    // &^
	Ampersand // &
	Pipe      // |
	Caret     // ^

	ShiftLeft  // <<
	ShiftRight // >>

	Equal          // ==
	NotEqual       // !=
	LessOrEqual    // <=
	GreaterOrEqual // >=
	Less           // <
	Greater        // >

	Exclamation // !

	Assign // =
	Define // :=

	Comma     // ,
	Period    // .
	Semicolon // ;
	Colons    // ::
	Colon     // :
	Hash      // #

	ParenLeft   // (
	BracketLeft // [
	BraceLeft   // {

	ParenRight   // )
	BracketRight // ]
	BraceRight   // }
)

func (k Kind) String() string {
	s := strings[k]
	if s == "" {
		return "<invalid token>"
	}
	return s
}

var strings = [...]string{
	Space:   "Space",
	Newline: "Newline",
	Comment: "Comment",

	Auto:     "auto",
	Break:    "break",
	Clone:    "clone",
	Continue: "continue",
	Else:     "else",
	False:    "false",
	For:      "for",
	If:       "if",
	Nil:      "nil",
	Return:   "return",
	True:     "true",

	Word: "Word",

	Integer:   "Integer",
	Character: "Character",
	String:    "String",

	Plus:     "+",
	Minus:    "-",
	Asterisk: "*",
	Slash:    "/",
	Percent:  "%",

	LogicalAnd: "&&",
	LogicalOr:  "||",

	AndNot:    "&^",
	Ampersand: "&",
	Pipe:      "|",
	Caret:     "^",

	ShiftLeft:  "<<",
	ShiftRight: ">>",

	Equal:          "==",
	NotEqual:       "!=",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",
	Less:           "<",
	Greater:        ">",

	Exclamation: "!",

	Assign: "=",
	Define: ":=",

	Comma:     ",",
	Period:    ".",
	Semicolon: ";",
	Colons:    "::",
	Colon:     ":",
	Hash:      "#",

	ParenLeft:   "(",
	BracketLeft: "[",
	BraceLeft:   "{",

	ParenRight:   ")",
	BracketRight: "]",
	BraceRight:   "}",
}
