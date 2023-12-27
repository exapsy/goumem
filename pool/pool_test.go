package pool

import (
	"testing"
	"unsafe"
)

func TestPool(t *testing.T) {
	t.Run("Allocating memory", func(t *testing.T) {
		pool, err := New(Options{
			Size: 15, // 15 bytes - account for alignment but not enough to store 3x32-bit integers
		})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		// Allocate 2 times 4 bytes (8 bytes total)
		addr, err := pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		addr.Set(456)
		if addr.Int() != 456 {
			t.Fatalf("Alloc() = %v, want 456", addr.Int())
		}

		addr2, err := pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		addr2.Set(456)
		if addr2.Int() != 456 {
			t.Fatalf("Alloc() = %v, want 456", addr2.Int())
		}

		// Free the pool
		err = pool.Close()
		if err != nil {
			t.Fatalf("Free() error = %v", err)
		}
	})

	t.Run("Freeing and reusing memory", func(t *testing.T) {
		var err error
		var addr *Ptr

		pool, err := New(Options{
			Size: 15, // 15 bytes - account for alignment but not enough to store 3x32-bit integers
		})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		// Allocate 2 times 4 bytes (8 bytes total)
		_, err = pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		addr, err = pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		// Free the first 4 bytes
		err = pool.Free(addr, 4)
		if err != nil {
			t.Fatalf("Free() error = %v", err)
		}

		// Allocate 4 bytes again (should work since we freed 4 bytes)
		addr3, err := pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		addr3.Set(456)
		if addr3.Int() != 456 {
			t.Errorf("Alloc() = %v, want 456", addr3.Int())
		}

		// Free the pool
		err = pool.Close()
		if err != nil {
			t.Fatalf("Free() error = %v", err)
		}
	})

	t.Run("Using the address memory directly", func(t *testing.T) {
		pool, err := New(Options{
			Size: 15, // 15 bytes - account for alignment but not enough to store 3x32-bit integers
		})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		// Allocate 4 bytes again (should work since we freed 4 bytes)
		addr4, err := pool.Alloc(4)
		if err != nil {
			t.Fatalf("Alloc() error = %v", err)
		}

		addr4.Set(456)
		if *(*uintptr)(unsafe.Pointer(addr4.Address())) != 456 {
			t.Fatalf("Alloc() = %v, want 456", addr4.Int())
		}

		// Free the pool
		err = pool.Close()
		if err != nil {
			t.Fatalf("Free() error = %v", err)
		}
	})
}
