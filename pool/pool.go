package pool

import (
	"fmt"
	goumem_syscall "github.com/exapsy/goumem/syscall"
	"unsafe"
)

type Pool struct {
	size        uint
	virtualAddr uintptr
	current     uintptr
	freed       []uintptr
}

type Options struct {
	Size uint
}

func New(opts Options) (*Pool, error) {
	addr, err := goumem_syscall.Alloc(uintptr(opts.Size))
	if err != nil {
		return nil, fmt.Errorf("failed to make MMAP syscall: %w", err)
	}

	return &Pool{
		size:        opts.Size,
		virtualAddr: addr,
		current:     addr,
		freed:       make([]uintptr, 0),
	}, nil
}

func (p Pool) PoolAddr() uintptr {
	return p.current - p.virtualAddr
}

// Alloc allocates memory inside the pool.
// Basically not really allocating memory, just returning the current address and incrementing it.
func (p *Pool) Alloc(size uint) (*Ptr, error) {
	for i, addr := range p.freed {
		if uintptr(size) <= addr {
			// Remove this block from the freed list
			p.freed = append(p.freed[:i], p.freed[i+1:]...)
			return &Ptr{
				VirtualAddr: addr,
				PoolAddr:    p.virtualAddr,
			}, nil
		}
	}
	if p.current+uintptr(size) > p.virtualAddr+uintptr(p.size) {
		return &Ptr{}, fmt.Errorf("pool is full")
	}

	addr := p.current
	p.current += uintptr(size)

	return &Ptr{
		VirtualAddr: addr,
		PoolAddr:    p.PoolAddr(),
	}, nil
}

// Free frees memory inside the pool.
// Basically not really freeing memory, just decrementing the current address.
func (p *Pool) Free(address *Ptr, size uint) error {
	addr := address.VirtualAddr
	if addr < p.virtualAddr || addr > p.virtualAddr+uintptr(p.size) {
		return fmt.Errorf("invalid virtualAddr")
	}

	p.freed = append(p.freed, addr)

	return nil
}

// Close frees the pool
func (p *Pool) Close() error {
	return goumem_syscall.Free(p.virtualAddr, uintptr(p.size))
}

type Ptr struct {
	VirtualAddr uintptr
	PoolAddr    uintptr
}

// Int returns the value of the pointer as an int
func (ptr Ptr) Int() int {
	return *(*int)(unsafe.Pointer(ptr.VirtualAddr))
}

// Ptr returns the value of the pointer as an uintptr
func (ptr Ptr) Uintptr() uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr.VirtualAddr))
}

// String returns the value of the pointer as a string
func (ptr Ptr) String() string {
	return *(*string)(unsafe.Pointer(ptr.VirtualAddr))
}

// Set sets the value of the pointer
func (ptr Ptr) Set(i interface{}) {
	switch v := i.(type) {
	case int:
		*(*int)(unsafe.Pointer(ptr.VirtualAddr)) = v
	case uintptr:
		*(*uintptr)(unsafe.Pointer(ptr.VirtualAddr)) = v
	case string:
		*(*string)(unsafe.Pointer(ptr.VirtualAddr)) = v
	}
}
