package goumem

import (
	"fmt"
	"sync"
	"unsafe"
)

type Ptr struct {
	virtualAddr uintptr
	poolAddr    uintptr
	mutex       sync.RWMutex
	size        uintptr
}

func (ptr *Ptr) Address() uintptr {
	return ptr.virtualAddr
}

// Int returns the value of the pointer as an int
func (ptr *Ptr) Int() int {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*int)(unsafe.Pointer(ptr.virtualAddr))
}

// Ptr returns the value of the pointer as an uintptr
func (ptr *Ptr) Uintptr() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*uintptr)(unsafe.Pointer(ptr.virtualAddr))
}

// String returns the value of the pointer as a string
func (ptr *Ptr) String() string {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*string)(unsafe.Pointer(ptr.virtualAddr))
}

// Set sets the value of the pointer
func (ptr *Ptr) Set(i interface{}) {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	switch v := i.(type) {
	case int:
		*(*int)(unsafe.Pointer(ptr.virtualAddr)) = v
	case uintptr:
		*(*uintptr)(unsafe.Pointer(ptr.virtualAddr)) = v
	case string:
		*(*string)(unsafe.Pointer(ptr.virtualAddr)) = v
	}
}

func (ptr *Ptr) Free() error {
	if ptr.size == 0 {
		return fmt.Errorf("cannot free a pointer without a size, it's probably allocated from a pool or an arena")
	}

	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.poolAddr, ptr.size)
}
