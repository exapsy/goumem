package goumem

import "testing"

func TestPointerInt(t *testing.T) {
	var i int
	ptr, err := NewInt(i)
	if err != nil {
		t.Fatal(err)
	}
	defer ptr.Free()

	if ptr.Value() != i {
		t.Fatal("Value() != i")
	}

	i = 1
	ptr.Set(i)
	if ptr.Value() != i {
		t.Fatal("Value() != i")
	}
}
