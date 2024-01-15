//go:build windows

package memsyscall

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type windowsSyscall struct{}

func New() *windowsSyscall {
	return &windowsSyscall{}
}

func (w *windowsSyscall) Alloc(size uintptr) (uintptr, error) {
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

func (w *windowsSyscall) Free(addr, size uintptr) error {
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

func (w *windowsSyscall) PageSize() uintptr {
	var info windows.SYSTEM_INFO
	windows.GetSystemInfo(&info)
	return uintptr(info.PageSize)
}
