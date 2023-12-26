package pool

import "testing"

func TestPool(t *testing.T) {
	pool, err := New(Options{
		Size: 1024,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	addr, err := pool.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc() error = %v", err)
	}

	addr.Set(456)
	if addr.Int() != 456 {
		t.Errorf("Alloc() = %v, want 456", addr.Int())
	}

	err = pool.Free(addr, 4)
	if err != nil {
		t.Fatalf("Free() error = %v", err)
	}
}
