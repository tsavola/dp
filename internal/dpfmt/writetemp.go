// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package dpfmt

import (
	"os"
	"path"
)

func ReplaceFile(filename string, data []byte) error {
	var (
		closed  bool
		renamed bool
	)

	f, err := os.CreateTemp(path.Dir(filename), ".*.dpfmt")
	if err != nil {
		return err
	}
	defer func() {
		if !renamed {
			os.Remove(f.Name())
		}
		if !closed {
			f.Close()
		}
	}()

	if _, err := f.Write(data); err != nil {
		return err
	}

	closed = true

	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), filename); err != nil {
		return err
	}

	renamed = true
	return nil
}
