package goumem

import (
	"sync"
	"unsafe"
)

type PointerInt8 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewInt8(i int8) (*PointerInt8, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerInt8{
		virtualAddr: virtualAddr,
	}
	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerInt8) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerInt8) Value() int8 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int8)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerInt8) Set(i int8) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*int8)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerInt8) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(int8(0)))
}
