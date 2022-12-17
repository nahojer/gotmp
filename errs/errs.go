// Package errs provides support for chaining errors.
package errs

import (
	"errors"
	"fmt"
)

// Chain represents an error chain.
type Chain []error

func (c Chain) Error() string {
	switch len(c) {
	case 0:
		return ""
	case 1:
		return c[0].Error()
	default:
		s := c[0].Error()
		for _, err := range c[1:] {
			s = fmt.Sprintf("%s: %s", s, err)
		}
		return s
	}
}

// Is reports whether any error in the chain matches target.
func (c Chain) Is(target error) bool {
	for _, err := range c {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As finds the first error in the chain that matches target, and if one is
// found, sets target to that error value and returns true. Otherwise, it
// returns false.
func (c Chain) As(target any) bool {
	for _, err := range c {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

// Unwrap returns itself.
func (c Chain) Unwrap() []error {
	return c
}
