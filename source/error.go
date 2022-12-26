// Copyright (c) 2023 Timo Savola
// SPDX-License-Identifier: BSD-3-Clause

package source

type wrappedError struct {
	msg string
	err error
}

func (e wrappedError) Error() string { return e.msg }
func (e wrappedError) Unwrap() error { return e.err }

// ErrorWithPositionPrefix enriches an error's message with position prefix.
// If the error object isn't position-aware and fallback is specified,
// "message" becomes "fallback: message".  Nil is passed through.
func ErrorWithPositionPrefix(err error, fallback string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(interface{ PositionError() string }); ok {
		return wrappedError{e.PositionError(), err}
	} else if fallback != "" {
		return wrappedError{fallback + ": " + err.Error(), err}
	}

	return err
}
