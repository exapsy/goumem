package goumem

import (
	"fmt"
	"syscall"
)

type memSyscallArgs struct {
	// The address of the memory to be allocated.
	addr uintptr
	// The size of the memory to be allocated.
	size uintptr
	// The protection of the memory to be allocated.
	prot uintptr
	// The flags of the memory to be allocated.
	flags uintptr
	// The file descriptor of the memory to be allocated.
	fd uintptr
	// The offset of the memory to be allocated.
	offset uintptr
}

func memSyscall(args memSyscallArgs) (uintptr, error) {
	syscallArgs := []uintptr{
		args.addr,
		args.size,
		args.prot,
		args.flags,
		args.fd,
		args.offset,
	}

	mem, _, errno := syscall.Syscall6(syscall.SYS_MMAP, syscallArgs[0], syscallArgs[1], syscallArgs[2], syscallArgs[3], syscallArgs[4], syscallArgs[5])
	if errno != 0 {
		return 0, fmt.Errorf("failed to make MMAP syscall: %w", errno)
	}

	return mem, nil
}

func freeSyscall(addr, size uintptr) error {
	syscallArgs := []uintptr{
		addr,
		size,
	}

	_, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, syscallArgs[0], syscallArgs[1], 0)
	if errno != 0 {
		return fmt.Errorf("failed to make MUNMAP syscall: %w", errno)
	}

	return nil
}
