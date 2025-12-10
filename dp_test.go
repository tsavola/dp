// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package dp_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/tsavola/dp/format"
	"github.com/tsavola/dp/lex"
	"github.com/tsavola/dp/parse"
	"github.com/tsavola/dp/source"

	. "github.com/tsavola/dp/internal/pan/mustcheck"
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

				parsed, err := parse.File(tokens)
				err = source.ErrorWithPositionPrefix(err, "")
				if err != nil {
					t.Fatalf("parse error:\n%v", err)
				}

				if false {
					for _, node := range parsed {
						t.Log(node)
					}
				}

				formatted := string(format.File(parsed))

				if false {
					t.Logf("formatted:\n%s", formatted)
				}
			})
		}
	}
}
