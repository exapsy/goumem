package goumem_syscall

import (
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

type MemoryAllocator interface {
	Mmap(size uintptr) (uintptr, error)
	Free(addr, size uintptr) error
}
