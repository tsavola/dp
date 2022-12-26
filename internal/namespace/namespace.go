// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package namespace

import (
	"strings"
)

func UnquoteImportPath(s string) (string, bool) {
	if !(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		panic(s)
	}
	s = s[1 : len(s)-1]
	for i, c := range s {
		if c == '\\' {
			// Return the valid prefix.
			return s[:i], false
		}
	}
	return s, true
}

func ImportPathNamespaces(path string) []string {
	if strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/") || strings.Contains(path, "//") {
		return nil
	}

	path = strings.ToLower(path)
	path = strings.Replace(path, "-", "_", -1)
	path = strings.Replace(path, ".", "_", -1)

	names := strings.Split(path, "/")
	names = append([]string{""}, names...) // For the :: prefix.

	var namespaces []string
	for i := range names {
		namespaces = append(namespaces, strings.Join(names[i:], "::"))
	}

	return namespaces
}
