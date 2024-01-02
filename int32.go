package goumem

import "unsafe"

type PointerInt32 struct {
	virtualAddr uintptr
}

func NewInt32(i int32) (*PointerInt32, error) {
	virtualAddr, err := mem.Alloc(unsafe.Sizeof(int32(0)))
	if err != nil {
		return nil, err
	}

	*(*int32)(unsafe.Pointer(virtualAddr)) = i

	return &PointerInt32{
		virtualAddr: virtualAddr,
	}, nil
}

func (ptr *PointerInt32) Address() uintptr {
	return ptr.virtualAddr
}

func (ptr *PointerInt32) Value() int32 {
	return *(*int32)(unsafe.Pointer(ptr.virtualAddr))
}

func (ptr *PointerInt32) Set(i int32) {
	*(*int32)(unsafe.Pointer(ptr.virtualAddr)) = i
}

func (ptr *PointerInt32) Free() error {
	return mem.Free(ptr.virtualAddr, unsafe.Sizeof(int32(0)))
}
