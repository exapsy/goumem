package goumem

import (
	"sync"
	"unsafe"
)

type PointerUint16 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewUint16(i uint16) (*PointerUint16, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerUint16{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerUint16) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerUint16) Value() uint16 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uint16)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerUint16) Set(i uint16) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*uint16)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerUint16) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(uint16(0)))
}
