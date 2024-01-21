package allocator

import (
	"fmt"
	memsyscall "github.com/exapsy/goumem/mem_syscall"
	"sync/atomic"
)

var (
	syscall                                     memsyscall.Syscall = memsyscall.New()
	PageSize                                                       = syscall.PageSize()
	allocThresholdWithoutAllocatingAnotherChunk uintptr            = PageSize / 2
)

type (
	MemoryAllocator interface {
		Alloc(size uintptr) (*AllocatedBlock, error)
		Free(block *AllocatedBlock) error
	}
	allocationStrategy interface {
		alloc(chunks *chunkList, size uintptr) (*AllocatedBlock, error)
		free(chunks *chunkList, block *AllocatedBlock) error
	}
	allocationPolicy interface {
		SelectChunk(chunks *chunkList, size uintptr) *chunk
	}
	chunkList struct {
		chunks *chunk
		len    int
	}
	chunk struct {
		size   atomic.Uintptr
		addr   uintptr
		blocks []*chunkBlock
		// freeBytes is meta-data for the total of free bytes left in the memory.
		freeBytes atomic.Uintptr
		next      *chunk
		prev      *chunk
	}
	chunkBlock struct {
		addr     atomic.Uintptr
		size     atomic.Uintptr
		isFree   atomic.Bool
		next     *chunkBlock
		prev     *chunkBlock
		nextFree *chunkBlock
		prevFree *chunkBlock
	}
	AllocatedBlock struct {
		size          uintptr
		addr          uintptr
		chunk         *chunk
		chunkBlockMem *chunkBlock
	}
)

func newChunkList() *chunkList {
	list := &chunkList{
		chunks: nil,
		len:    1,
	}

	newChunk, err := newChunk(PageSize, nil, nil)
	if err != nil {
		panic(err)
	}

	list.chunks = newChunk

	return list
}

func (cl *chunkList) alloc(size uintptr) (block *AllocatedBlock, chunkIndex int, blockIndex int, err error) {
	head := cl.chunks
	for head != nil {
		var addr uintptr
		addr, blockIndex, err = head.getFreeAddrForSize(size)
		if err != nil {
			return nil, 0, 0, err
		}

		return &AllocatedBlock{
			addr:          addr,
			size:          size,
			chunk:         head,
			chunkBlockMem: head.blocks[blockIndex],
		}, chunkIndex, blockIndex, nil
	}
	return nil, 0, 0, nil
}

func (cl *chunkList) freeBytes() uintptr {
	var freeBytes uintptr
	head := cl.chunks
	for head != nil {
		head.freeBytes.Add(freeBytes)
		head = head.next
	}

	return freeBytes
}

func (cl *chunkList) freeChunk(chunkMem *chunk) error {
	// remove chunk from list
	if chunkMem.prev != nil { // not the first chunk
		chunkMem.prev.next = chunkMem.next
	}

	if chunkMem.next != nil { // not the last chunk
		chunkMem.next.prev = chunkMem.prev
	}

	// free memory
	err := syscall.Free(chunkMem.addr, chunkMem.size.Load())
	if err != nil {
		return fmt.Errorf("could not free memory: %w", err)
	}

	cl.len--

	return nil
}

func newChunk(size uintptr, previous, next *chunk) (*chunk, error) {
	var addr uintptr
	var err error

	memoryAlignedSize := size
	if size%PageSize != 0 {
		memoryAlignedSize = size + (PageSize - size%PageSize) // align to next page
	}

	// alloc space from kernel
	addr, err = syscall.Alloc(memoryAlignedSize)
	if err != nil {
		return nil, fmt.Errorf("could not alloc memory: %w", err)
	}

	c := &chunk{
		addr: addr,
		prev: previous,
		next: next,
		blocks: []*chunkBlock{
			{
				next:     nil,
				prev:     nil,
				nextFree: nil,
				prevFree: nil,
			},
		},
	}

	c.size.Store(memoryAlignedSize)
	c.freeBytes.Store(memoryAlignedSize)
	c.blocks[0].isFree.Store(true)
	c.blocks[0].size.Store(memoryAlignedSize)
	c.blocks[0].addr.Store(addr)

	return c, nil
}

func (c *chunk) get(index int) *chunk {
	head := c
	for i := 0; i < index; i++ {
		head = head.next
	}

	return head
}

// allocAndAppendNewChunk allocates a new chunk at the end of the tail of chunks
// of size being the page-size of the system.
func (c *chunk) allocAndAppendNewChunk() (chunck *chunk, err error) {
	var addr uintptr

	// alloc space from kernel
	addr, err = syscall.Alloc(PageSize)
	if err != nil {
		return nil, fmt.Errorf("could not alloc memory: %w", err)
	}

	// Get the tail
	head := c
	for head.next != nil {
		head = head.next
	}

	// Append new chunk on the tail
	head.next = &chunk{
		addr: addr,
		blocks: []*chunkBlock{
			{},
		},
		next: nil,
		prev: head,
	}

	head.freeBytes.Add(PageSize)
	head.size.Add(PageSize)
	head.blocks[0].isFree.Store(true)
	head.blocks[0].size.Store(PageSize)
	head.blocks[0].addr.Store(addr)
	head.blocks[0].isFree.Store(true)

	return head.next, nil
}

func (c *chunk) splitAndGetFirstPart(block *chunkBlock, size uintptr) (uintptr, error) {
	// handle blocks
	firstBlock := block
	firstBlock.isFree.Store(false)
	secondBlock := &chunkBlock{
		prev:     firstBlock,
		next:     firstBlock.next,
		nextFree: firstBlock.nextFree,
		prevFree: firstBlock.prevFree,
	}
	secondBlock.addr.Store(firstBlock.addr.Load() + size + 1)
	secondBlock.size.Store(firstBlock.size.Load() - size)
	secondBlock.isFree.Store(true)
	firstBlock.next = secondBlock

	if firstBlock.prevFree != nil {
		firstBlock.prevFree.nextFree = secondBlock
	}
	if firstBlock.nextFree != nil {
		firstBlock.nextFree.nextFree = secondBlock.nextFree
	}

	// handle chunk
	c.freeBytes.Add(-size)
	c.blocks = append(c.blocks, secondBlock)

	return firstBlock.addr.Load(), nil
}

// mergeAdjacent merges every adjacent block nearby.
// Reduces fragmentation of memory, even when seemed not necessary.
// Used by [free] method.
func (c *chunkBlock) mergeAdjacent() {
	// check for backwards adjacent free blocks
	head := c
	for head.isFree.Load() && head.addr.Load()+c.size.Load()+1 == c.addr.Load() {
		head.size.Add(c.size.Load())
		head.prev = c.next
	}

	// check for forward adjacent free blocks
	head = c
	for head.isFree.Load() && c.addr.Load()+c.size.Load() == head.addr.Load()-1 {
		c.size.Add(head.size.Load())
		c.next = head.next
	}

	if c.prevFree != nil {
		prevNextFree := c.prevFree.nextFree
		c.prevFree.nextFree = c

		c.prevFree = prevNextFree
	}

	if c.nextFree != nil {
		c.nextFree.prevFree = c
		nextPrevFree := c.nextFree.prevFree
		c.nextFree = nextPrevFree
	}
}

func (c *chunk) getFreeAddrForSize(requestedSize uintptr) (addr uintptr, blockIndex int, err error) {
	// O(n) max cycles
	// 2 branches - 1 with 2 other branches
	for i, block := range c.blocks {
		// split block if it's free and has enough space and occupy the first part of it
		if block.isFree.Load() && block.size.Load() > requestedSize {
			addr, err := c.splitAndGetFirstPart(block, requestedSize)
			if err != nil {
				return 0, 0, err
			}

			return addr, i, nil
		} else if block.isFree.Load() {
			// free blocks should be already merged by free method
			// so, we do not need to check for adjacent free blocks
			// it should already be at its maximum size

			if block.size.Load() > requestedSize {
				// check again if it's the requested size and split if
				addr, err := c.splitAndGetFirstPart(block, requestedSize)
				if err != nil {
					return 0, 0, err
				}

				return addr, i, nil
			} else if block.size.Load() == requestedSize {
				// alloc whole block
				addr := block.addr
				c.freeBytes.Add(-requestedSize)
				block.isFree.Store(false)
				return addr.Load(), i, nil
			}
		}
	}

	return 0, 0, fmt.Errorf("could not alloc memory in chunk")
}

func (b *AllocatedBlock) Addr() uintptr {
	return b.addr
}

func (b *AllocatedBlock) Size() uintptr {
	return b.size
}
