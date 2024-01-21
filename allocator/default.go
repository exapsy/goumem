package allocator

import "fmt"

var (
	Default = func() MemoryAllocator { return NewDefaultMemoryAllocator() }
)

type defaultAllocPolicy struct{}

func newDefaultAllocPolicy() allocationPolicy {
	return &defaultAllocPolicy{}
}

// SelectChunk selects a chunk from the list of chunks with more or equal free bytes than the size.
func (p *defaultAllocPolicy) SelectChunk(chunkList *chunkList, size uintptr) *chunk {
	if allocThresholdWithoutAllocatingAnotherChunk >= size {
		// threshold not reached
		for i := 0; i < chunkList.len; i++ {
			// select first chunk with enough free bytes
			chunk := chunkList.chunks
			if chunk != nil && chunk.freeBytes.Load() >= size {
				return chunk
			}
		}
	}

	// chunk with this amount of free bytes not found
	// or threshold is reached
	// allocate new chunk
	lastChunk := chunkList.chunks.get(chunkList.len - 1)
	c, err := newChunk(size, lastChunk, nil)
	if err != nil {
		fmt.Println("error allocating new chunk: ", err)
		return nil
	}

	lastChunk.next = c
	c.prev = lastChunk
	chunkList.len++

	return c
}

type defaultAllocStrategy struct {
	allocationPolicy allocationPolicy
}

func newDefaultAllocStrategy(allocPolicy allocationPolicy) allocationStrategy {
	return &defaultAllocStrategy{
		allocationPolicy: allocPolicy,
	}
}

func (s *defaultAllocStrategy) alloc(chunks *chunkList, size uintptr) (*AllocatedBlock, error) {
	c := s.allocationPolicy.SelectChunk(chunks, size)
	if c == nil {
		return nil, fmt.Errorf("no chunk found")
	}

	for i := 0; i < len(c.blocks); i++ {
		block := c.blocks[i]
		if block.isFree.Load() && block.size.Load() >= size {
			addr, err := c.splitAndGetFirstPart(block, size)
			if err != nil {
				return nil, err
			}

			return &AllocatedBlock{
				size:          size,
				addr:          addr,
				chunk:         c,
				chunkBlockMem: block,
			}, nil
		}
	}

	return nil, fmt.Errorf(
		`no free block found. This should not happen
as the policy should have selected a chunk with free bytes. chunk: %+v`, c,
	)
}

func (s *defaultAllocStrategy) free(chunks *chunkList, block *AllocatedBlock) error {

	block.chunkBlockMem.isFree.Store(true)

	// merge adjacent blocks
	// That is, if there are any adjacent blocks that are free.
	// Improves fragmentation of memory,
	// and we don't have to search
	// and merge everytime on allocation very fragmented memory.
	block.chunkBlockMem.mergeAdjacent()

	block.chunk.freeBytes.Add(block.size)

	if block.chunk.freeBytes == block.chunk.size &&
		chunks.freeBytes() > PageSize {

		if block.chunk.prev == nil { // first chunk
			block.chunk.next.prev = nil
		}

		if block.chunk.next == nil { // last chunk
			block.chunk.prev.next = nil
		}

		// free chunk
		// if it is the last chunk
		// and the previous chunk has enough free bytes
		err := chunks.freeChunk(block.chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

type defaultMemoryAllocator struct {
	strategy allocationStrategy
	chunks   *chunkList
}

func NewDefaultMemoryAllocator() MemoryAllocator {
	cl := newChunkList()
	return &defaultMemoryAllocator{
		strategy: newDefaultAllocStrategy(newDefaultAllocPolicy()),
		chunks:   cl,
	}
}

func (a *defaultMemoryAllocator) Alloc(size uintptr) (*AllocatedBlock, error) {
	return a.strategy.alloc(a.chunks, size)
}

func (a *defaultMemoryAllocator) Free(block *AllocatedBlock) error {
	return a.strategy.free(a.chunks, block)
}
