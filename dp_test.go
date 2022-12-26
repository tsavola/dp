// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package dp_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/tsavola/dp/lex"
	"github.com/tsavola/dp/source"

	. "import.name/pan/mustcheck"
)

func Test(t *testing.T) {
	for _, e := range Must(os.ReadDir("testdata")) {
		if name, ok := strings.CutSuffix(e.Name(), "_test.dp"); ok {
			filename := path.Join("testdata", e.Name())

			t.Run(name, func(t *testing.T) {
				input := string(Must(os.ReadFile(filename)))

				tokens, err := lex.File(source.Location(filename), input)
				err = source.ErrorWithPositionPrefix(err, "")
				if err != nil {
					t.Fatalf("tokenization error:\n%v", err)
				}

				if false {
					s := "tokens:\n"
					for _, tok := range tokens {
						s += "·"
						s += tok.Source
					}
					t.Log(s)
				}
			})
		}
	}
}
