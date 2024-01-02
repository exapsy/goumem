package goumem

import (
	"sync"
	"unsafe"
)

type PointerUint8 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewUint8(i uint8) (*PointerUint8, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerUint8{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerUint8) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerUint8) Value() uint8 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uint8)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerUint8) Set(i uint8) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*uint8)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerUint8) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(uint8(0)))
}
