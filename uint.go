package goumem

import (
	"sync"
	"unsafe"
)

type PointerUint struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewUint(i uint) (*PointerUint, error) {
	ptr := &PointerUint{}

	var err error
	ptr.virtualAddr, err = mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerUint) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerUint) Value() uint {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uint)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerUint) Set(i uint) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*uint)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerUint) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(uint(0)))
}
