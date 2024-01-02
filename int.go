package goumem

import (
	"sync"
	"unsafe"
)

type PointerInt struct {
	virtualAddr uintptr
	mutex       sync.RWMutex
}

func NewInt(i int) (*PointerInt, error) {
	ptr := &PointerInt{}

	var err error
	ptr.virtualAddr, err = mem.Alloc(unsafe.Sizeof(i))
	if err != nil {
		return nil, err
	}

	ptr.Set(i)

	return ptr, nil
}

func (ptr *PointerInt) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerInt) Value() int {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerInt) Set(i int) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	*(*int)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerInt) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(int(0)))
}
