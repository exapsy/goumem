package goumem

import (
	"sync"
	"unsafe"
)

type PointerInt64 struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewInt64(i int64) (*PointerInt64, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr := &PointerInt64{
		virtualAddr: virtualAddr,
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerInt64) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerInt64) Value() int64 {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int64)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerInt64) Set(i int64) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*int64)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerInt64) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(int64(0)))
}
