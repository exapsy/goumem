//go:build windows

// Package goumem_syscall
package allocator

import (
	"fmt"
	"golang.org/x/sys/windows"
)

func New() MemoryAllocator {
	return Windows{}
}

type Windows struct{}

func (w Windows) Alloc(size uintptr) (uintptr, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")

	r1, _, err := virtualAlloc.Call(
		0,
		size,
		uintptr(0x1000), // MEM_COMMIT
		uintptr(0x04),   // PAGE_READWRITE
	)
	if r1 == 0 {
		return 0, fmt.Errorf("failed to make VirtualAlloc allocator: %w", err)
	}

	return r1, nil
}

func (w Windows) Free(addr, size uintptr) error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	virtualFree := kernel32.NewProc("VirtualFree")

	r1, _, err := virtualFree.Call(
		addr,
		0,
		uintptr(0x8000), // MEM_RELEASE
	)
	if r1 == 0 {
		return fmt.Errorf("failed to make VirtualFree allocator: %w", err)
	}

	return nil
}
