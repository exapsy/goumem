package allocator

import (
	"fmt"
	memsyscall "github.com/exapsy/goumem/mem_syscall"
	"sync"
)

var (
	syscall                                     memsyscall.Syscall = memsyscall.New()
	PageSize                                                       = syscall.PageSize()
	AllocThresholdWithoutAllocatingAnotherChunk uintptr            = PageSize / 2
)

type (
	MemoryAllocator interface {
		Alloc(size uintptr) (*AllocatedBlock, error)
		Free(block *AllocatedBlock) error
	}
	AllocationStrategy interface {
		Allocate(chunks *chunkList, size uintptr) (*AllocatedBlock, error)
		Free(chunks *chunkList, block *AllocatedBlock) error
	}
	AllocationPolicy interface {
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
		addr   uintptr
		size   uintptr
		isFree bool
		next   *chunkBlock
		prev   *chunkBlock
	}
	AllocatedBlock struct {
		size       uintptr
		addr       uintptr
		chunkIndex int
		blockIndex int
	}
)

func newChunkList() *chunkList {
	list := chunkList{
		len: 1,
	}

	newChunk, err := newChunk(nil, nil)
	if err != nil {
		panic(err)
	}

	list.chunks = newChunk

	return nil
}

func (cl *chunkList) alloc(size uintptr) (block *AllocatedBlock, chunkIndex int, blockIndex int, err error) {
	cl.chunks.alloc(size)
	return nil, 0, 0, nil
}

func newChunk(previous, next *chunk) (*chunk, error) {
	var addr uintptr
	var err error

	// Allocate space from kernel
	addr, err = syscall.Alloc(PageSize)
	if err != nil {
		return nil, fmt.Errorf("could not allocate memory: %w", err)
	}

	return &chunk{
		addr:      addr,
		size:      PageSize,
		freeBytes: PageSize,
		prev:      previous,
		next:      next,
		blocks: []*chunkBlock{
			{
				addr:   addr,
				size:   PageSize,
				isFree: true,
				next:   nil,
				prev:   nil,
			},
		},
	}, nil
}

func (c *chunk) alloc(size uintptr) (addr uintptr, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// check if there is enough space in the chunk
	if c.freeBytes < size {
		return 0, fmt.Errorf("not enough space in chunk")
	}

	// get the first free address for the requested size
	addr, _, err = c.getFreeAddrForSize(size)
	if err != nil {
		return 0, fmt.Errorf("could not get free address for size: %w", err)
	}

	return addr, nil
}

// allocNewChunk allocates a new chunk at the end of the tail of chunks
// of size being the page-size of the system.
func (c *chunk) allocNewChunk() (chunck *chunk, err error) {
	var addr uintptr

	// Allocate space from kernel
	addr, err = syscall.Alloc(PageSize)
	if err != nil {
		return nil, fmt.Errorf("could not allocate memory: %w", err)
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

func (c *chunk) splitAndGetFirst(block *chunkBlock, size uintptr) (uintptr, error) {
	c.freeBytes -= size
	firstBlock := block
	firstBlock.isFree = false
	secondBlock := &chunkBlock{
		addr:   firstBlock.addr + size + 1,
		size:   firstBlock.size - size,
		isFree: true,
		prev:   firstBlock,
		next:   firstBlock.next,
	}
	firstBlock.next = secondBlock
	return firstBlock.addr, nil
}

// mergeAdjacent merges every adjacent block nearby.
// Reduces fragmentation of memory, even when seemed not necessary.
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
}

func (c *chunk) getFreeAddrForSize(requestedSize uintptr) (addr uintptr, blockIndex int, err error) {
	// check free blocks and adjust them if needed and have adjacent ones
	// Adjust as in example: (x = occupied block(s), - is free block(s) )
	// [xx | -  | - | x | - | xxxx | - ] ->
	// [xx | -- | x | - | xxxx | - ] merged adjacent free block
	//
	// O(n) max cycles
	// 2 branches - 1 with 2 other branches
	for i, block := range c.blocks {
		// split block if it's free and has enough space and occupy the first part of it
		if block.isFree && block.size > requestedSize {
			addr, err := c.splitAndGetFirst(block, requestedSize)
			if err != nil {
				return 0, 0, err
			}

			return addr, i, nil
		} else if block.isFree {
			// check if block is free &
			// merge all adjacent blocks
			block.mergeAdjacent()

			if block.size > requestedSize {
				// check again if it's the requested size and split if
				addr, err := c.splitAndGetFirst(block, requestedSize)
				if err != nil {
					return 0, 0, err
				}

				return addr, i, nil
			} else if block.size == requestedSize {
				// allocate whole block
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
	//newChunk, err := c.allocNewChunk()
	//if err != nil {
	//	return 0, 0, fmt.Errorf("could not allocate new chunk")
	//}

	// get the first chunk which has the requested requestedSize
	//first, err := newChunk.splitAndGetFirst(newChunk.blocks[0], requestedSize)
	//if err != nil {
	//	return 0, 0, err
	//}

	return 0, 0, fmt.Errorf("could not allocate memory in chunk")
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

func (b *AllocatedBlock) Addr() uintptr {
	return b.addr
}

func (b *AllocatedBlock) Size() uintptr {
	return b.size
}
