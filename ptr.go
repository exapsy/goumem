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

// Address returns the address of the pointer
// You can use it to directly read/write to the memory,
// but you have to be careful about the type of the pointer.
//
// It is the most unsafe method of this package for writing or reading from or to the memory.
// But it's also the most volatile and pools don't consider types, just allocations.
// So it's the most efficient method for such jobs.
//
// If you're sure about the type, you're free to use methods of this type such as Int(), String(), etc.
// But if you're not sure about the type, consider just dereferencing the pointer and using the value.
func (ptr *Ptr) Address() uintptr {
	return ptr.virtualAddr
}

// PtrTypeInt returns the value of the pointer as an int
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

// Byte returns the value of the pointer as a byte
func (ptr *Ptr) Byte() byte {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*byte)(unsafe.Pointer(ptr.virtualAddr))
}

// Set sets the value of the pointer
func (ptr *Ptr) Set(i interface{}) error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	switch v := i.(type) {
	case int:
		*(*int)(unsafe.Pointer(ptr.virtualAddr)) = v
	case uintptr:
		*(*uintptr)(unsafe.Pointer(ptr.virtualAddr)) = v
	case string:
		*(*string)(unsafe.Pointer(ptr.virtualAddr)) = v
	case byte:
		*(*byte)(unsafe.Pointer(ptr.virtualAddr)) = v
	default:
		*(*interface{})(unsafe.Pointer(ptr.virtualAddr)) = v
	}

	return nil
}

func (ptr *Ptr) Free() error {
	if ptr.size == 0 {
		return fmt.Errorf("cannot free a pointer without a size, it's probably allocated from a pool or an arena")
	}

	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.poolAddr, ptr.size)
}
