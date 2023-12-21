package goumem

import (
	"syscall"
	"unsafe"
)

type Int uintptr

func NewInt() (Int, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(int(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Int(u), nil
}

func (ui Int) Set(i int) {
	*(*int)(unsafe.Pointer(ui)) = i
}

func (ui Int) Val() int {
	return *(*int)(unsafe.Pointer(ui))
}

func (ui Int) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(int(0)))
}

type Uintptr uintptr

func NewUintptr() (Uintptr, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uintptr(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uintptr(u), nil
}

func (ui Uintptr) Set(i uintptr) {
	*(*uintptr)(unsafe.Pointer(ui)) = i
}

func (ui Uintptr) Val() uintptr {
	return *(*uintptr)(unsafe.Pointer(ui))
}

func (ui Uintptr) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uintptr(0)))
}

type Uint8 uintptr

func NewUint8() (Uint8, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uint8(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uint8(u), nil
}

func (ui Uint8) Set(i uint8) {
	*(*uint8)(unsafe.Pointer(ui)) = i
}

func (ui Uint8) Val() uint8 {
	return *(*uint8)(unsafe.Pointer(ui))
}

func (ui Uint8) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uint8(0)))
}

type Uint16 uintptr

func NewUint16() (Uint16, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uint16(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uint16(u), nil
}

func (ui Uint16) Set(i uint16) {
	*(*uint16)(unsafe.Pointer(ui)) = i
}

func (ui Uint16) Val() uint16 {
	return *(*uint16)(unsafe.Pointer(ui))
}

func (ui Uint16) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uint16(0)))
}

type Uint32 uintptr

func NewUint32() (Uint32, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uint32(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uint32(u), nil
}

func (ui Uint32) Set(i uint32) {
	*(*uint32)(unsafe.Pointer(ui)) = i
}

func (ui Uint32) Val() uint32 {
	return *(*uint32)(unsafe.Pointer(ui))
}

func (ui Uint32) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uint32(0)))
}

type Uint64 uintptr

func NewUint64() (Uint64, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uint64(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uint64(u), nil
}

func (ui Uint64) Set(i uint64) {
	*(*uint64)(unsafe.Pointer(ui)) = i
}

func (ui Uint64) Val() uint64 {
	return *(*uint64)(unsafe.Pointer(ui))
}

func (ui Uint64) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uint64(0)))
}

type Uint uintptr

func NewUint() (Uint, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(uint(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Uint(u), nil
}

func (ui Uint) Set(i uint) {
	*(*uint)(unsafe.Pointer(ui)) = i
}

func (ui Uint) Val() uint {
	return *(*uint)(unsafe.Pointer(ui))
}

func (ui Uint) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(uint(0)))
}

type String uintptr

func NewString() (String, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   255,
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return String(u), nil
}

func (ui String) Set(s string) {
	*(*string)(unsafe.Pointer(ui)) = s
}

func (ui String) Val() string {
	return *(*string)(unsafe.Pointer(ui))
}

func (ui String) Free() error {
	return freeSyscall(uintptr(ui), 255)
}

type Byte uintptr

func NewByte() (Byte, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(byte(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Byte(u), nil
}

func (ui Byte) Set(b byte) {
	*(*byte)(unsafe.Pointer(ui)) = b
}

func (ui Byte) Val() byte {
	return *(*byte)(unsafe.Pointer(ui))
}

func (ui Byte) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(byte(0)))
}

type Int64 uintptr

func NewInt64() (Int64, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(int64(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})
	if err != nil {
		return 0, err
	}

	return Int64(u), nil
}

func (ui Int64) Set(i int64) {
	*(*int64)(unsafe.Pointer(ui)) = i
}

func (ui Int64) Val() int64 {
	return *(*int64)(unsafe.Pointer(ui))
}

func (ui Int64) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(int64(0)))
}

type Int32 uintptr

func NewInt32() (Int32, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(int32(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0,
	})

	if err != nil {
		return 0, err
	}

	return Int32(u), nil
}

func (ui Int32) Set(i int32) {
	*(*int32)(unsafe.Pointer(ui)) = i
}

func (ui Int32) Val() int32 {
	return *(*int32)(unsafe.Pointer(ui))
}

func (ui Int32) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(int32(0)))
}

type Int16 uintptr

func NewInt16() (Int16, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(int16(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0})

	if err != nil {
		return 0, err
	}

	return Int16(u), nil
}

func (ui Int16) Set(i int16) {
	*(*int16)(unsafe.Pointer(ui)) = i
}

func (ui Int16) Val() int16 {
	return *(*int16)(unsafe.Pointer(ui))
}

func (ui Int16) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(int16(0)))
}

type Int8 uintptr

func NewInt8() (Int8, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(int8(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0})

	if err != nil {
		return 0, err
	}

	return Int8(u), nil
}

func (ui Int8) Set(i int8) {
	*(*int8)(unsafe.Pointer(ui)) = i
}

func (ui Int8) Val() int8 {
	return *(*int8)(unsafe.Pointer(ui))
}

func (ui Int8) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(int8(0)))
}

type Float64 uintptr

func NewFloat64() (Float64, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(float64(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0})

	if err != nil {
		return 0, err
	}

	return Float64(u), nil
}

func (ui Float64) Set(i float64) {
	*(*float64)(unsafe.Pointer(ui)) = i
}

func (ui Float64) Val() float64 {
	return *(*float64)(unsafe.Pointer(ui))
}

func (ui Float64) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(float64(0)))
}

type Float32 uintptr

func NewFloat32() (Float32, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(float32(0)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0})

	if err != nil {
		return 0, err
	}

	return Float32(u), nil
}

func (ui Float32) Set(i float32) {
	*(*float32)(unsafe.Pointer(ui)) = i
}

func (ui Float32) Val() float32 {
	return *(*float32)(unsafe.Pointer(ui))
}

func (ui Float32) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(float32(0)))
}

type Bool uintptr

func NewBool() (Bool, error) {
	u, err := memSyscall(memSyscallArgs{
		addr:   0,
		size:   unsafe.Sizeof(bool(false)),
		prot:   syscall.PROT_READ | syscall.PROT_WRITE,
		flags:  syscall.MAP_ANONYMOUS | syscall.MAP_PRIVATE,
		fd:     0,
		offset: 0})

	if err != nil {
		return 0, err
	}

	return Bool(u), nil
}

func (ui Bool) Set(i bool) {
	*(*bool)(unsafe.Pointer(ui)) = i
}

func (ui Bool) Val() bool {
	return *(*bool)(unsafe.Pointer(ui))
}

func (ui Bool) Free() error {
	return freeSyscall(uintptr(ui), unsafe.Sizeof(bool(false)))
}
