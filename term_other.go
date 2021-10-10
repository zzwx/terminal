//go:build !windows
// +build !windows

// Package terminal
package terminal

import (
	"os"
)

func EnableVirtualTerminalProcessing(f *os.File, enable bool) error {
	return nil
}
