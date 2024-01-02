package goumem

import (
	"sync"
	"unsafe"
)

type PointerUint32 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewUint32(i uint32) (*PointerUint32, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerUint32{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerUint32) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerUint32) Value() uint32 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uint32)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerUint32) Set(i uint32) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*uint32)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerUint32) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(uint32(0)))
}
