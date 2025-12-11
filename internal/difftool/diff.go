// Copyright (c) 2024 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package difftool

import (
	"os"
	"os/exec"

	. "github.com/tsavola/dp/internal/pan/mustcheck"
)

var (
	command = "diff"
	options = "-u"
)

func MustDiff(filename string, data []byte) {
	cmd := exec.Command(command, options, filename, "/dev/stdin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	w := Must(cmd.StdinPipe())
	Check(cmd.Start())
	Must(w.Write(data))
	Check(w.Close())
	cmd.Wait() // diff exit status is nonzero if there are differences.
}
