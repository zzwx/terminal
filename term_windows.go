//go:build windows
// +build windows

package terminal

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32Dll    = syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode = kernel32Dll.NewProc("SetConsoleMode")
	getConsoleMode = kernel32Dll.NewProc("GetConsoleMode")
)

func EnableVirtualTerminalProcessing(f *os.File, enable bool) error {
	var mode uint32
	h := f.Fd()
	_, _, err := getConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode)))
	if err != nil {
		if e := err.(syscall.Errno); e != 0 {
			return err
		}
	}
	if enable {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	} else {
		mode &^= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	}
	_, _, err = setConsoleMode.Call(h, uintptr(mode))
	if err != nil {
		if e := err.(syscall.Errno); e != 0 {
			return err
		}
	}
	return nil
}
