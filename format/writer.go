// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package format

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type writer struct {
	*bytes.Buffer
}

func (w writer) currentLineLen() int {
	s := w.String()
	i := strings.LastIndex(s, "\n")
	if i >= 0 {
		s = s[i+1:]
	}
	return utf8.RuneCountInString(s)
}
