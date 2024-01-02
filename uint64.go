package goumem

import (
	"sync"
	"unsafe"
)

type PointerUint64 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewUint64(i uint64) (*PointerUint64, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerUint64{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

// Address returns the virtual address of the pointer.
// You should never use this to write to the pointer.
//
// It is not thread-safe, unless you are only reading.
// If you are only reading, you should use Value() instead.
func (ptr *PointerUint64) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerUint64) Value() uint64 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uint64)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerUint64) Set(i uint64) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*uint64)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerUint64) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(uint64(0)))
}
