// Package memsyscall are calls directly to the system, without any middleware,
// directly to and from the system.
//
// It makes sure the correct calls are made for the correct system/OS.
package memsyscall

type Syscall interface {
	Alloc(size uintptr) (addr uintptr, err error)
	Free(addr uintptr, size uintptr) (err error)
	PageSize() (size uintptr)
}
