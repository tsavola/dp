// Copyright (c) 2022 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package parse

import (
	"errors"

	"github.com/tsavola/dp/internal/position"
	"github.com/tsavola/dp/source"
	"github.com/tsavola/dp/token"
	"import.name/pan"
)

type poser interface {
	pos() source.Position
}

func parse[ScanState poser, Result any](s ScanState, parsers ...func(ScanState) (ScanState, Result)) (ScanState, Result) {
	var errs []error

	for _, f := range parsers {
		var after ScanState
		var node Result

		err := pan.Recover(func() {
			after, node = f(s)
		})
		if err == nil {
			return after, node
		}

		errs = append(errs, err)
	}

	if len(errs) == 1 {
		var e position.Error
		if errors.As(errs[0], &e) {
			pan.Panic(errs[0])
		}
	}

	panic(pan.Wrap(position.NewError(s.pos(), "syntax error", errs...)))
}

func parseNakedList[T comparable](s scan, assignment bool, parsers ...func(scan) (scan, T)) (scan, []T) {
	var results []T

	for {
		switch s.peek().Kind {
		case token.BraceLeft, token.BraceRight, token.Colon, token.Comment, token.Define, token.Newline, token.Semicolon:
			return s, results

		case token.Assign:
			if assignment {
				return s, results
			}
		}

		var r T
		s, r = parse(s, parsers...)

		var zero T
		if r != zero {
			results = append(results, r)
		}
	}
}

func parseListUntil[T comparable](s scan, stop func(scan) (scan, bool), parsers ...func(scan) (scan, T)) (scan, []T) {
	var results []T

	for {
		if end, ok := stop(s); ok {
			return end, results
		}

		var r T
		s, r = parse(s, parsers...)

		var zero T
		if r != zero {
			results = append(results, r)
		}
	}
}
