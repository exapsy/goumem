//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris

package memsyscall

import (
	"fmt"
	"syscall"
)

type unixSyscall struct{}

func New() *unixSyscall {
	return &unixSyscall{}
}

func (u *unixSyscall) Alloc(size uintptr) (uintptr, error) {
	mem, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,
		size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANONYMOUS|syscall.MAP_PRIVATE,
		0,
		0,
	)
	if errno != 0 {
		return 0, fmt.Errorf("failed to make MMAP allocator: %w", errno)
	}

	return mem, nil
}

func (u *unixSyscall) Free(addr uintptr, size uintptr) error {
	syscallArgs := []uintptr{
		addr,
		size,
	}

	_, _, errno := syscall.Syscall(
		syscall.SYS_MUNMAP,
		syscallArgs[0],
		syscallArgs[1],
		0,
	)
	if errno != 0 {
		return fmt.Errorf("failed to make MUNMAP allocator: %w", errno)
	}

	return nil
}

func (u *unixSyscall) PageSize() uintptr {
	return uintptr(syscall.Getpagesize())
}
