package allocator

import (
	"fmt"
	memsyscall "github.com/exapsy/goumem/mem_syscall"
	"sync"
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
		size   uintptr
		addr   uintptr
		blocks []*chunkBlock
		// freeBytes is meta-data for the total of free bytes left in the memory.
		freeBytes uintptr
		mutex     sync.Mutex
		next      *chunk
		prev      *chunk
	}
	chunkBlock struct {
		addr     uintptr
		size     uintptr
		isFree   bool
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
		freeBytes += head.freeBytes
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
	err := syscall.Free(chunkMem.addr, chunkMem.size)
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

	return &chunk{
		addr:      addr,
		size:      memoryAlignedSize,
		freeBytes: memoryAlignedSize,
		prev:      previous,
		next:      next,
		blocks: []*chunkBlock{
			{
				addr:     addr,
				size:     memoryAlignedSize,
				isFree:   true,
				next:     nil,
				prev:     nil,
				nextFree: nil,
				prevFree: nil,
			},
		},
	}, nil
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
		size:      PageSize,
		freeBytes: PageSize,
		addr:      addr,
		blocks: []*chunkBlock{
			{
				size: PageSize,
				addr: addr,
			},
		},
		mutex: sync.Mutex{},
		next:  nil,
		prev:  head,
	}

	return head.next, nil
}

func (c *chunk) splitAndGetFirstPart(block *chunkBlock, size uintptr) (uintptr, error) {
	// handle blocks
	firstBlock := block
	firstBlock.isFree = false
	secondBlock := &chunkBlock{
		addr:     firstBlock.addr + size + 1,
		size:     firstBlock.size - size,
		isFree:   true,
		prev:     firstBlock,
		next:     firstBlock.next,
		nextFree: firstBlock.nextFree,
		prevFree: firstBlock.prevFree,
	}
	firstBlock.next = secondBlock

	if firstBlock.prevFree != nil {
		firstBlock.prevFree.nextFree = secondBlock
	}
	if firstBlock.nextFree != nil {
		firstBlock.nextFree.nextFree = secondBlock.nextFree
	}

	// handle chunk
	c.freeBytes -= size
	c.blocks = append(c.blocks, secondBlock)

	return firstBlock.addr, nil
}

// mergeAdjacent merges every adjacent block nearby.
// Reduces fragmentation of memory, even when seemed not necessary.
// Used by [free] method.
func (c *chunkBlock) mergeAdjacent() {
	// check for backwards adjacent free blocks
	head := c
	for head.isFree && head.addr+c.size+1 == c.addr {
		head.size += c.size
		head.prev = c.next
	}

	// check for forward adjacent free blocks
	head = c
	for head.isFree && c.addr+c.size == head.addr-1 {
		c.size += head.size
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
		if block.isFree && block.size > requestedSize {
			addr, err := c.splitAndGetFirstPart(block, requestedSize)
			if err != nil {
				return 0, 0, err
			}

			return addr, i, nil
		} else if block.isFree {
			// free blocks should be already merged by free method
			// so, we do not need to check for adjacent free blocks
			// it should already be at its maximum size

			if block.size > requestedSize {
				// check again if it's the requested size and split if
				addr, err := c.splitAndGetFirstPart(block, requestedSize)
				if err != nil {
					return 0, 0, err
				}

				return addr, i, nil
			} else if block.size == requestedSize {
				// alloc whole block
				addr := block.addr
				c.freeBytes -= requestedSize
				block.isFree = false
				return addr, i, nil
			}
		}
	}

	// TODO: This solution below is very localized.
	//      It does not allow for having the index of the chunk
	// 		Thus why the chunkList was made,
	//		and this kind of logic should be moved there,
	// 		so, we can know in which chunk the allocation happened and transfer this metadata
	//      to the respective AllocatedBlock.

	// create new chunk
	//newChunk, err := c.allocAndAppendNewChunk()
	//if err != nil {
	//	return 0, 0, fmt.Errorf("could not alloc new chunk")
	//}

	// get the first chunk which has the requested requestedSize
	//first, err := newChunk.splitAndGetFirstPart(newChunk.blocks[0], requestedSize)
	//if err != nil {
	//	return 0, 0, err
	//}

	return 0, 0, fmt.Errorf("could not alloc memory in chunk")
}

func (c *chunkBlock) append(block *chunkBlock) error {
	if c == nil {
		return fmt.Errorf("chunkBlock is nil")
	}

	head := c
	for head.next != nil {
		head = head.next
	}

	head.next = block

	return nil
}

func (c *chunkBlock) get(index int) *chunkBlock {
	head := c
	for i := 0; i < index; i++ {
		head = head.next
	}

	return head
}

func (b *AllocatedBlock) Addr() uintptr {
	return b.addr
}

func (b *AllocatedBlock) Size() uintptr {
	return b.size
}
