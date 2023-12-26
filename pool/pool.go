package pool

import (
	"fmt"
	"github.com/exapsy/goumem/allocator"
	"sync"
	"unsafe"
)

var (
	ErrPoolFull = fmt.Errorf("pool is full")
)

var (
	mem = allocator.New()
)

type Pool struct {
	size        uint
	virtualAddr uintptr
	current     uintptr
	freed       []uintptr
	mutex       sync.Mutex
}

type Options struct {
	Size uint
}

func New(opts Options) (*Pool, error) {
	allocSize := opts.Size
	addr, err := mem.Alloc(uintptr(allocSize))
	if err != nil {
		return nil, fmt.Errorf("failed to make MMAP allocator: %w", err)
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
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Calculate the alignment for the requested size
	alignment := uintptr(size)
	if alignment < unsafe.Alignof(uintptr(0)) {
		alignment = unsafe.Alignof(uintptr(0))
	}

	// Align the current address
	misalignment := p.current % alignment
	adjustment := uintptr(0)
	if misalignment != 0 {
		adjustment = alignment - misalignment
	}

	if p.current+adjustment+uintptr(size) > p.virtualAddr+uintptr(p.size) {
		return &Ptr{}, fmt.Errorf("pool is full")
	}

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

	addr := p.current
	p.current += uintptr(size)

	return &Ptr{
		VirtualAddr: addr,
		PoolAddr:    p.PoolAddr(),
	}, nil
}

// Free frees memory inside the pool.
// Basically not really freeing memory, just decrementing the current address
func (p *Pool) Free(address *Ptr, size uint) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	addr := address.VirtualAddr
	if addr < p.virtualAddr || addr > p.virtualAddr+uintptr(size) {
		return fmt.Errorf("invalid virtualAddr")
	}

	p.freed = append(p.freed, addr)

	return nil
}

// Close frees the pool
func (p *Pool) Close() error {
	return mem.Free(p.virtualAddr, uintptr(p.size))
}

type Ptr struct {
	VirtualAddr uintptr
	PoolAddr    uintptr
	mutex       sync.RWMutex
}

// Int returns the value of the pointer as an int
func (ptr *Ptr) Int() int {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int)(unsafe.Pointer(ptr.VirtualAddr))
}

// Ptr returns the value of the pointer as an uintptr
func (ptr *Ptr) Uintptr() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uintptr)(unsafe.Pointer(ptr.VirtualAddr))
}

// String returns the value of the pointer as a string
func (ptr *Ptr) String() string {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*string)(unsafe.Pointer(ptr.VirtualAddr))
}

// Set sets the value of the pointer
func (ptr *Ptr) Set(i interface{}) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	switch v := i.(type) {
	case int:
		*(*int)(unsafe.Pointer(ptr.VirtualAddr)) = v
	case uintptr:
		*(*uintptr)(unsafe.Pointer(ptr.VirtualAddr)) = v
	case string:
		*(*string)(unsafe.Pointer(ptr.VirtualAddr)) = v
	}
}
