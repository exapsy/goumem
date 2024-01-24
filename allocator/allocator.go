package allocator

import (
	"fmt"
	memsyscall "github.com/exapsy/goumem/mem_syscall"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var (
	syscall                                     memsyscall.Syscall = memsyscall.New()
	PageSize                                                       = syscall.PageSize()
	allocThresholdWithoutAllocatingAnotherChunk uintptr            = PageSize / 2
)

var (
	ErrAllocatedBlockAlreadyFreed  = fmt.Errorf("goumem: allocated block already freed")
	ErrAllocatedBlockDifferentSize = fmt.Errorf("goumem: allocated block different size")
)

type (
	MemoryAllocator interface {
		Alloc(size uintptr) (*AllocatedBlock, error)
		Free(block *AllocatedBlock) error
		Copy(dst, src *AllocatedBlock) error
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
		flags         AllocatedBlockFlags
	}
	AllocatedBlockFlags uintptr
)

const (
	AllocatedBlockFlagsNone AllocatedBlockFlags = 0
	AllocatedBlockFlagsFree AllocatedBlockFlags = 1 << iota
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

func (b *AllocatedBlock) IsFreed() bool {
	return b.flags&AllocatedBlockFlagsFree != 0
}

func setTypes(ptr unsafe.Pointer, val reflect.Value) {
	field := val
	switch field.Kind() {
	case reflect.Int:
		*(*int)(ptr) = int(field.Int())
	case reflect.Int8:
		*(*int8)(ptr) = int8(field.Int())
	case reflect.Int16:
		*(*int16)(ptr) = int16(field.Int())
	case reflect.Int32:
		*(*int32)(ptr) = int32(field.Int())
	case reflect.Int64:
		*(*int64)(ptr) = field.Int()
	case reflect.Float64:
		*(*float64)(ptr) = field.Float()
	case reflect.String:
		*(*string)(ptr) = field.String()
	case reflect.Bool:
		*(*bool)(ptr) = field.Bool()
	case reflect.Uintptr:
		*(*uintptr)(ptr) = uintptr(field.Uint())
	case reflect.Uint:
		*(*uint)(ptr) = uint(field.Uint())
	case reflect.Uint8:
		*(*uint8)(ptr) = uint8(field.Uint())
	case reflect.Uint16:
		*(*uint16)(ptr) = uint16(field.Uint())
	case reflect.Uint32:
		*(*uint32)(ptr) = uint32(field.Uint())
	case reflect.Uint64:
		*(*uint64)(ptr) = uint64(field.Uint())
	case reflect.Complex64:
		*(*complex64)(ptr) = (complex64)(field.Complex())
	case reflect.Complex128:
		*(*complex128)(ptr) = field.Complex()
	case reflect.Ptr:
		setPtr(ptr, field)
	case reflect.Struct:
		setStruct(ptr, field)
	case reflect.Array, reflect.Slice:
		setArrayOrSlice(ptr, field)
	case reflect.Map:
		setMap(ptr, val)
	default:
		panic("unsupported type")
	}
}

func setStruct(ptr unsafe.Pointer, val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		default:
			setTypes(ptr, field)
		}
	}
}

// getANyType
//
// @atPtr: pointer to the address of the value to be set
//
// @from: the value to be set, from which also the type is inferred. So, make sure the type of the value is correct, with respect to the type of the pointer.
func getAnyTypeTo(atPtr unsafe.Pointer, from reflect.Value) {
	var realType reflect.Type
	var realValueRefl reflect.Value
	refRealValueRefl := from

	realValueRefl = from
	realType = from.Type()

	// figure out if its a pointer
	// if it is, get the element type
	switch realType.Kind() {
	case reflect.Ptr:
		realType = from.Type().Elem()
	default:
		realType = from.Type()
	}

	switch realType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		refRealValueRefl.Elem().SetInt(int64(*(*int)(atPtr)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		refRealValueRefl.Elem().SetUint(uint64(*(*uint)(atPtr)))
	case reflect.Float32, reflect.Float64:
		refRealValueRefl.Elem().SetFloat(float64(*(*float64)(atPtr)))
	case reflect.Complex64, reflect.Complex128:
		refRealValueRefl.Elem().SetComplex(complex128(*(*complex128)(atPtr)))
	case reflect.Array, reflect.Slice:
		getArrayOrSlice(atPtr, realValueRefl)
	case reflect.Struct:
		getStruct(atPtr, realValueRefl)
	case reflect.Ptr:
		getPtr(atPtr, realValueRefl)
	case reflect.Map:
		getMap(atPtr, realValueRefl)
	case reflect.String:
		refRealValueRefl.Elem().SetString(*(*string)(atPtr))
	case reflect.Bool:
		refRealValueRefl.Elem().SetBool(*(*bool)(atPtr))
	default:
		panic("unsupported type")
	}
}

func getMap(ptr unsafe.Pointer, val reflect.Value) {
	for _, key := range val.MapKeys() {
		elem := val.MapIndex(key)
		switch elem.Kind() {
		default:
			getAnyTypeTo(ptr, elem)
		}
	}
}

func getPtr(ptr unsafe.Pointer, val reflect.Value) {
	switch val.Kind() {
	default:
		// get elem
		val = val.Elem()
		getAnyTypeTo(ptr, val)
	}
}

func getStruct(ptr unsafe.Pointer, val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		valAtIndex := val.Elem().Index(i)
		ptrAtIndex := unsafe.Pointer(uintptr(ptr) + val.Elem().Index(i).Type().Size())
		getAnyTypeTo(ptrAtIndex, valAtIndex)
	}
}

func setArrayOrSlice(ptr unsafe.Pointer, val reflect.Value) {
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		switch elem.Kind() {
		default:
			getAnyTypeTo(ptr, val)
		}
	}
}

func getArrayOrSlice(ptr unsafe.Pointer, val reflect.Value) {
	for i := 0; i < val.Len(); i++ {
		valAtIndex := val.Elem().Index(i)
		ptr = unsafe.Pointer(uintptr(ptr) + val.Elem().Index(i).Type().Size())
		getAnyTypeTo(ptr, valAtIndex)
	}
}

func setMap(ptr unsafe.Pointer, val reflect.Value) {
	for _, key := range val.MapKeys() {
		elem := val.MapIndex(key)
		switch elem.Kind() {
		default:
			setTypes(ptr, elem)
		}
	}
}

func setPtr(ptr unsafe.Pointer, val reflect.Value) {
	switch val.Kind() {
	default:
		// get elem
		val = val.Elem()
		setTypes(ptr, val)
	}
}

func (cb *chunkBlock) copy(dst *chunkBlock) error {
	src := cb
	srcSize := src.size.Load()
	destSize := dst.size.Load()

	if srcSize != destSize {
		return ErrAllocatedBlockDifferentSize
	}

	srcSlice := (*(*[1 << 30]byte)(unsafe.Pointer(src.addr.Load())))[:srcSize:srcSize]
	destSlice := (*(*[1 << 30]byte)(unsafe.Pointer(dst.addr.Load())))[:destSize:destSize]

	copy(destSlice, srcSlice)

	return nil
}

func Get[T any](block *AllocatedBlock) T {
	return *(*T)(unsafe.Pointer(block.Addr()))
}

func Set[T any](block *AllocatedBlock, value T) {
	*(*T)(unsafe.Pointer(block.Addr())) = value
}
