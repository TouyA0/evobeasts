//go:build windows

package main

import (
	"os"
	"syscall"
)

var procFlushConsoleInputBuffer = syscall.NewLazyDLL("kernel32.dll").NewProc("FlushConsoleInputBuffer")

// flushInput jette les touches tapées pendant les animations,
// pour qu'un Entrée en trop ne valide pas le menu suivant tout seul.
func flushInput() {
	pendingKeys = nil
	procFlushConsoleInputBuffer.Call(os.Stdin.Fd())
}
