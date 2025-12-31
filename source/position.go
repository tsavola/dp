// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Package source describes source code locations.
package source

type Position struct {
	Path       string
	Line       int // 1-based
	Column     int // 1-based
	ByteOffset int // 0-based
}

func Location(path string) Position {
	return Position{
		Path:   path,
		Line:   1,
		Column: 1,
	}
}
