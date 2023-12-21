package goumem

import (
	"testing"
	"unsafe"
)

func TestInt(t *testing.T) {
	u, err := NewInt()
	if err != nil {
		t.Errorf("Int() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Int() = %v, want 456", *(*int)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestUintptr(t *testing.T) {
	u, err := NewUintptr()
	if err != nil {
		t.Errorf("Uintptr() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Uintptr() = %v, want 456", *(*uintptr)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestString(t *testing.T) {
	u, err := NewString()
	if err != nil {
		t.Errorf("String() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set("456")
	if u.Val() != "456" {
		t.Errorf("String() = %v, want 456", *(*string)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestUint(t *testing.T) {
	u, err := NewUint()
	if err != nil {
		t.Errorf("Uint() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Uint() = %v, want 456", *(*uint)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestInt64(t *testing.T) {
	u, err := NewInt64()
	if err != nil {
		t.Errorf("Int64() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Int64() = %v, want 456", *(*int64)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestUint64(t *testing.T) {
	u, err := NewUint64()
	if err != nil {
		t.Errorf("Uint64() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Uint64() = %v, want 456", *(*uint64)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestFloat64(t *testing.T) {
	u, err := NewFloat64()
	if err != nil {
		t.Errorf("Float64() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Float64() = %v, want 456", *(*float64)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestFloat32(t *testing.T) {
	u, err := NewFloat32()
	if err != nil {
		t.Errorf("Float32() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(456)
	if u.Val() != 456 {
		t.Errorf("Float32() = %v, want 456", *(*float32)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestBool(t *testing.T) {
	u, err := NewBool()
	if err != nil {
		t.Errorf("Bool() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(true)
	if u.Val() != true {
		t.Errorf("Bool() = %v, want true", *(*bool)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}

func TestByte(t *testing.T) {
	u, err := NewByte()
	if err != nil {
		t.Errorf("Byte() error = %v", err)
		return
	}

	// Test that the memory is writable
	u.Set(1)
	if u.Val() != 1 {
		t.Errorf("Byte() = %v, want 1", *(*byte)(unsafe.Pointer(u)))
	}

	err = u.Free()
	if err != nil {
		t.Errorf("Free() error = %v", err)
		return
	}
}
