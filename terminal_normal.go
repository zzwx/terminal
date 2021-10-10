//go:build !js
// +build !js

package terminal

import (
	"golang.org/x/term"
)

// IsTerminal returns whether the fd passed in is a terminal or not
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// ReadPassword reads a line of input from a terminal without local echo. This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) {
	return term.ReadPassword(fd)
}
