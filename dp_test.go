// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package dp_test

import (
	"os"
	"path"
	"strings"
	"testing"

	. "import.name/pan/mustcheck"
)

func Test(t *testing.T) {
	for _, e := range Must(os.ReadDir("testdata")) {
		if name, ok := strings.CutSuffix(e.Name(), "_test.dp"); ok {
			filename := path.Join("testdata", e.Name())

			t.Run(name, func(t *testing.T) {
				input := string(Must(os.ReadFile(filename)))
			})
		}
	}
}
