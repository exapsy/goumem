# Goumem - Manual Memory Allocation Library

Who said Go doesn't have manual memory allocation?

**NO NEED TO CALL C** and context switch to manually allocate memory.

Goumem is a library that provides manual memory allocation functions for Go,
the way C does it, by using `mmap`.

## Installation

```bash
go get github.com/exapsy/goumem
```

## Example - Create a pool

```go
package main

import (
	"fmt"
	goumempool "github.com/exapsy/goumem/pool"
)

func main() {
	pool, err := goumempool.New(Options{
		// Size of the pool in bytes
        Size: 15,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Allocate 2 times 4 bytes (8 bytes total)
	addr, err := pool.Alloc(4)
	if err != nil {
		panic(fmt.Errorf("Alloc() error = %v", err))
	}

	addr.Set(456)
	if addr.Int() != 456 {
		panic(fmt.Errorf("Alloc() = %v, want 456", addr.Int()))
	}

	addr2, err := pool.Alloc(4)
	if err != nil {
		panic(fmt.Errorf("Alloc() error = %v", err))
	}

	addr2.Set(456)
	if addr2.Int() != 456 {
		panic(fmt.Errorf("Alloc() = %v, want 456", addr2.Int()))
	}

	// Free the first 4 bytes
	err = pool.Free(addr, 4)
	if err != nil {
		panic(fmt.Errorf("Free() error = %v", err))
	}

	// Allocate 4 bytes again (should work since we freed 4 bytes)
	addr3, err := pool.Alloc(4)
	if err != nil {
		panic(fmt.Errorf("Alloc() error = %v", err))
	}

	addr3.Set(456)
	if addr3.Int() != 456 {
		panic(fmt.Errorf("Alloc() = %v, want 456", addr3.Int()))
	}

	// Free the pool
	err = pool.Close()
	if err != nil {
		panic(fmt.Errorf("Free() error = %v", err))
	}
}
```
