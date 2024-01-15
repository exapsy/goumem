package allocator

var (
	Default MemoryAllocator = &defaultAllocator{
		chunkList: newChunkList(),
	}
)

type defaultAllocator struct {
	chunks *chunkList
}

func (a *defaultAllocator) Alloc(size uintptr) (*AllocatedBlock, error) {
	alloc, blockIndex, err := a.chunks.alloc(size)
	if err != nil {
		return AllocatedBlock{}, err
	}

	return AllocatedBlock{
		addr:       alloc,
		size:       size,
		blockIndex: blockIndex,
	}, nil
}

func (a *defaultAllocator) Free(b *AllocatedBlock) error {
	return nil
}
