// Copyright (c) 2025 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package dp

import (
	"testing"

	"github.com/tsavola/dp/format"
	"github.com/tsavola/dp/lex"
	"github.com/tsavola/dp/parse"
	"github.com/tsavola/dp/source"
)

func Fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		tokens, err := lex.File(source.Location(t.Name()), input)
		if err != nil {
			return
		}

		parsed, err := parse.File(tokens)
		if err != nil {
			return
		}

		formatted := string(format.File(parsed))
		if formatted == input {
			return
		}

		tokens, err = lex.File(source.Location("formatted-"+t.Name()), formatted)
		if err != nil {
			t.Fatal("formatted lex error:", err)
		}

		parsed, err = parse.File(tokens)
		if err != nil {
			t.Fatal("formatted parse error:", err)
		}

		reformatted := string(format.File(parsed))
		if reformatted != formatted {
			t.Fatal("format is not stable")
		}
	})
}
