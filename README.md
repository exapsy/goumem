# Goumem - Manual Memory Allocation Library

Who said Go doesn't have manual memory allocation?

**NO NEED TO CALL C** and context switch to manually allocate memory.

Goumem is a library that provides manual memory allocation functions for Go,
the way C does it, by using `mmap`.

## Installation

```bash
go get github.com/exapsy/goumem
```

## Example - Create an integer

```go
package main

import (
	"fmt"
	"github.com/exapsy/goumem"
)

func main() { 
    // Allocate an integer 
    i, err := goumem.NewInt()
    if err != nil {
        panic(err)
    }
		
    // Set the value of the integer 
    i.Set(42)
    
    // Get the value of the integer 
    intVal := i.Val()
    fmt.Println(intVal) // 42
    
    // Free the memory allocated for the integer 
    err = i.Free()
    if err != nil {
        panic(fmt.Errorf("error freeing memory: %v", err))
    }
}
```

## Example - Create a string

```go
package main

import (
    "fmt"
    "github.com/exapsy/goumem"
)

func main() {
    // Allocate a string
    s, err := goumem.NewString()
    if err != nil {
        panic(err)
    }
    
    // Set the value of the string
    s.Set("Hello World!")
    
    // Get the value of the string
    strVal := s.Val()
    fmt.Println(strVal) // Hello World!
    
    // Free the memory allocated for the string
    err = s.Free()
    if err != nil {
        panic(fmt.Errorf("error freeing memory: %v", err))
    }
}
```