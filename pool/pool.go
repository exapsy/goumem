package pool

import (
	"fmt"
	goumem_syscall "github.com/exapsy/goumem/syscall"
	"unsafe"
)

type Pool struct {
	size    uint
	address uintptr
	current uintptr
}

type Options struct {
	Size uint
}

func New(opts Options) (*Pool, error) {
	addr, err := goumem_syscall.Mmap(uintptr(opts.Size))
	if err != nil {
		return nil, fmt.Errorf("failed to make MMAP syscall: %w", err)
	}

	return &Pool{
		size:    opts.Size,
		address: addr,
		current: addr,
	}, nil
}

// Alloc allocates memory inside the pool.
// Basically not really allocating memory, just returning the current address and incrementing it.
func (p *Pool) Alloc(size uint) (Uintptr, error) {
	if p.current+uintptr(size) > p.address+uintptr(p.size) {
		return 0, fmt.Errorf("pool is full")
	}

	addr := p.current
	p.current += uintptr(size)

	return Uintptr(addr), nil
}

// Free frees memory inside the pool.
// Basically not really freeing memory, just decrementing the current address.
func (p *Pool) Free(address Uintptr, size uint) error {
	addr := uintptr(address)
	if addr < p.address || addr > p.address+uintptr(p.size) {
		return fmt.Errorf("invalid address")
	}

	if addr+uintptr(size) > p.current {
		return fmt.Errorf("invalid size")
	}

	p.current -= uintptr(size)

	return nil
}

type Uintptr uintptr

// Int returns the value of the pointer as an int
func (ui Uintptr) Int() int {
	return *(*int)(unsafe.Pointer(ui))
}

// Uintptr returns the value of the pointer as an uintptr
func (ui Uintptr) Uintptr() uintptr {
	return *(*uintptr)(unsafe.Pointer(ui))
}

// String returns the value of the pointer as a string
func (ui Uintptr) String() string {
	return *(*string)(unsafe.Pointer(ui))
}

// Set sets the value of the pointer
func (ui Uintptr) Set(i interface{}) {
	switch v := i.(type) {
	case int:
		*(*int)(unsafe.Pointer(ui)) = v
	case uintptr:
		*(*uintptr)(unsafe.Pointer(ui)) = v
	case string:
		*(*string)(unsafe.Pointer(ui)) = v
	}
}

// Close frees the pool
func (p *Pool) Close() error {
	return goumem_syscall.Free(p.address, uintptr(p.size))
}
