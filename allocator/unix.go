//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris

// Package goumem_syscall
package allocator

import (
	"fmt"
	"syscall"
)

type Protection uintptr

const (
	ProtectionRead  Protection = syscall.PROT_READ
	ProtectionWrite Protection = syscall.PROT_WRITE
)

type Flag uintptr

const (
	FlagAnonymous Flag = syscall.MAP_ANONYMOUS
	FlagPrivate   Flag = syscall.MAP_PRIVATE
)

func New() MemoryAllocator {
	return Unix{}
}

type Unix struct{}

func (mau Unix) Alloc(size uintptr) (uintptr, error) {
	syscallArgs := []uintptr{
		0,
		size,
		uintptr(ProtectionRead) | uintptr(ProtectionWrite),
		uintptr(FlagAnonymous) | uintptr(FlagPrivate),
		0,
		0,
	}

	mem, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		syscallArgs[0],
		syscallArgs[1],
		syscallArgs[2],
		syscallArgs[3],
		syscallArgs[4],
		syscallArgs[5],
	)
	if errno != 0 {
		return 0, fmt.Errorf("failed to make MMAP allocator: %w", errno)
	}

	return mem, nil
}

func (mau Unix) Free(addr, size uintptr) error {
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
