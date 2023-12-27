package allocator

type MemoryAllocator interface {
	Alloc(size uintptr) (uintptr, error)
	Free(addr, size uintptr) error
}
