package goumem

import (
	"sync"
	"unsafe"
)

type PointerInt16 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewInt16(i int16) (*PointerInt16, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerInt16{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerInt16) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerInt16) Value() int16 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int16)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerInt16) Set(i int16) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*int16)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerInt16) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(int16(0)))
}
