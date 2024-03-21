// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

// Command dpfmt is a source code formatter.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/tsavola/dp/ast"
	"github.com/tsavola/dp/format"
	"github.com/tsavola/dp/internal/dpfmt"
	"github.com/tsavola/dp/internal/revise"
	"github.com/tsavola/dp/lex"
	"github.com/tsavola/dp/parse"
	"github.com/tsavola/dp/source"
	"import.name/pan"

	. "import.name/pan/mustcheck"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file>\n", os.Args[0])
		flag.PrintDefaults()
	}

	var (
		old   = flag.Bool("old", false, "parse old language version")
		diff  = flag.Bool("d", false, "display diffs instead of rewriting files")
		write = flag.Bool("w", false, "write result to (source) file instead of stdout")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	filename := flag.Arg(0)

	err := pan.Recover(func() {
		program(filename, *old, *diff, *write)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, source.ErrorWithPositionPrefix(err, filename))
		os.Exit(1)
	}
}

func program(filename string, old, diff, write bool) {
	pos := source.Location(filename)
	input := string(Must(os.ReadFile(filename)))

	var parsed []ast.FileChild
	if !old {
		parsed = Must(parse.File(Must(lex.File(pos, input))))
	} else {
		parsed = Must(revise.File(pos, input))
	}

	output := format.File(parsed)

	if diff {
		cmd := exec.Command("diff", "-u", filename, "/dev/stdin")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		w := Must(cmd.StdinPipe())
		Check(cmd.Start())
		Must(w.Write(output))
		Check(w.Close())
		cmd.Wait() // diff exit status is nonzero if there are differences.
	} else if !write {
		Must(os.Stdout.Write(output))
		return
	}

	if write {
		Check(dpfmt.ReplaceFile(filename, output))
	}
}
