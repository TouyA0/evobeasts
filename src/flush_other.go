//go:build !windows

package main

func flushInput() {
	pendingKeys = nil
}