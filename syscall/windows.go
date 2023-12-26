//go:build windows

// Package goumem_syscall
package goumem_syscall

import (
	"fmt"
	"golang.org/x/sys/windows"
)

func Free(addr, size uintptr) error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	virtualFree := kernel32.NewProc("VirtualFree")

	r1, _, err := virtualFree.Call(
		addr,
		0,
		uintptr(0x8000), // MEM_RELEASE
	)
	if r1 == 0 {
		return fmt.Errorf("failed to make VirtualFree syscall: %w", err)
	}

	return nil
}
