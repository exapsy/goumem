package goumem

import (
	"fmt"
	"sync"
	"unsafe"
)

func ErrStringSizeTooBig(originalSize uintptr, newStringSize uintptr) error {
	return errStringSizeTooBig{
		originalSize:  originalSize,
		newStringSize: newStringSize,
	}
}

type errStringSizeTooBig struct {
	originalSize  uintptr
	newStringSize uintptr
}

func (e errStringSizeTooBig) Error() string {
	return fmt.Sprintf("string size is too big: original size %d bytes > %d bytes", e.originalSize, e.newStringSize)
}

type PointerStr struct {
	virtualAddr uintptr
	poolAddr    uintptr
	mutex       sync.RWMutex
	size        uintptr
}

func NewString(s string) (*PointerStr, error) {
	poolAddr, err := mem.Alloc(uintptr(len(s)))
	if err != nil {
		return nil, err
	}

	ptr := &PointerStr{
		poolAddr: poolAddr,
		size:     uintptr(len(s)),
	}

	ptr.virtualAddr = poolAddr

	ptr.Set(s)

	return ptr, nil
}

func (ptr *PointerStr) Address() uintptr {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return ptr.virtualAddr
}

func (ptr *PointerStr) Value() string {
	ptr.mutex.RLock()
	defer ptr.mutex.RUnlock()

	return *(*string)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerStr) Set(s string) error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	if ptr.size < uintptr(len(s)) {
		return ErrStringSizeTooBig(ptr.size, uintptr(len(s)))
	}

	*(*string)(unsafe.Pointer(ptr.virtualAddr)) = s

	return nil
}

func (ptr *PointerStr) Free() error {
	ptr.mutex.Lock()
	defer ptr.mutex.Unlock()

	return mem.Free(ptr.poolAddr, ptr.size)
}
