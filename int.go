package goumem

import (
	"github.com/exapsy/goumem/allocator"
	"sync"
	"unsafe"
)

type PointerInt struct {
	allocatedBlock *allocator.AllocatedBlock
	mutex          sync.RWMutex
}

func NewInt(i int) (*PointerInt, error) {
	ptr := &PointerInt{}

	var err error
	var block *allocator.AllocatedBlock
	block, err = mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr.allocatedBlock = block

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerInt) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.allocatedBlock.Addr()
}

func (ptr *PointerInt) Value() int {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int)(unsafe.Pointer(ptr.allocatedBlock.Addr()))
}

func (ptr *PointerInt) Set(i int) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*int)(unsafe.Pointer(ptr.allocatedBlock.Addr())) = i
}

func (ptr *PointerInt) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.allocatedBlock)
}

func init() {
	if mem == nil {
		mem = allocator.Default()
	}
}
