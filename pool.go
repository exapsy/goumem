package goumem

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
	freed       []FreedBlock
	mutex       sync.Mutex
}

type FreedBlock struct {
	Addr uintptr
	Size uint
}

type PoolOptions struct {
	Size uint
}

func NewPool(opts PoolOptions) (*Pool, error) {
	allocSize := opts.Size
	addr, err := mem.Alloc(uintptr(allocSize))
	if err != nil {
		return nil, fmt.Errorf("failed to make MMAP allocator: %w", err)
	}

	return &Pool{
		size:        opts.Size,
		virtualAddr: addr,
		current:     addr,
		freed:       make([]FreedBlock, 0),
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
		return nil, ErrPoolFull
	}

	for i, block := range p.freed {
		if size <= block.Size {
			// Remove this block from the freed list
			p.freed = append(p.freed[:i], p.freed[i+1:]...)
			return &Ptr{
				virtualAddr: block.Addr,
				poolAddr:    p.virtualAddr,
			}, nil
		}
	}

	addr := p.current
	p.current += uintptr(size)

	return &Ptr{
		virtualAddr: addr,
		poolAddr:    p.PoolAddr(),
	}, nil
}

// Free frees memory inside the pool.
// Basically not really freeing memory, just decrementing the current address
func (p *Pool) Free(address *Ptr, size uint) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	addr := address.virtualAddr
	if addr < p.virtualAddr || addr > p.virtualAddr+uintptr(size) {
		return fmt.Errorf("invalid virtualAddr")
	}

	// Search for a freed block that is adjacent to this one
	// If any extend this block to include the freed block
	// This prevents the freed list from growing too large -
	// or in other words, the over-fragmentation of the freed list
	for i, block := range p.freed {
		if block.Addr+uintptr(block.Size) == addr {
			// Extend this block to include the freed block
			p.freed[i].Size += size
			return nil
		} else if addr+uintptr(size) == block.Addr {
			// Extend the freed block to include this block
			p.freed[i].Addr -= uintptr(size)
			p.freed[i].Size += size
			return nil
		}
	}

	// No adjacent blocks so just append this one to the freed list
	p.freed = append(p.freed, FreedBlock{
		Addr: addr,
		Size: size,
	})

	return nil
}

// Close frees the pool
func (p *Pool) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return mem.Free(p.virtualAddr, uintptr(p.size))
}
