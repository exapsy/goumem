# Goumem - Manual Memory Allocation Library

Who said Go doesn't have manual memory allocation?

Goumem is a library that provides manual memory allocation functions for Go.

## Installation

```bash
go get github.com/exapsy/goumem
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/exapsy/goumem"
)

func main() { 
    // Allocate an integer 
    i, p := goumem.Int(2) // Returns an int and a pointer to it 
    fmt.Printf("i: %d, p: %p\n", i, p)
	
    // Write to the allocated memory 
    *p = 3
    fmt.Printf("i: %d, p: %p\n", i, p)
	
    // Free the allocated memory
    goumem.FreeInt(p)
}
```
