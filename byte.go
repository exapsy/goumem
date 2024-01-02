package goumem

import "unsafe"

type PointerByte struct {
	virtualAddr uintptr
}

func (ptr *PointerByte) Address() uintptr {
	return ptr.virtualAddr
}

func (ptr *PointerByte) Value() byte {
	return *(*byte)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerByte) Set(i byte) {
	*(*byte)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerByte) Free() error {
	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(byte(0)))
}
