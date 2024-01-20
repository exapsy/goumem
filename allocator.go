package goumem

import "github.com/exapsy/goumem/allocator"

var (
	mem allocator.MemoryAllocator
)

func SetMemoryAllocator(m allocator.MemoryAllocator) {
	mem = m
}
