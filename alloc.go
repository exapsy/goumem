package goumem

import (
	"fmt"
	"github.com/exapsy/goumem/allocator"
)

var (
	mem = allocator.New()
)

func Alloc(size uintptr) (*Ptr, error) {
	addr, err := mem.Alloc(size)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory: %w", err)
	}

	return &Ptr{
		virtualAddr: addr,
		poolAddr:    0,
		size:        size,
	}, nil
}
