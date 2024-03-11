# Goumem - Manual Memory Allocation Library

<p align="center">
  <img src="./docs/assets/logo.png" alt="goumem" width="350"/>
</p>

Who said Go doesn't have manual memory allocation?

**NO NEED TO CALL C** and context switch to manually allocate memory.

Goumem is a library that provides manual memory allocation functions for Go,
the way C does it, by using `mmap` for unix systems or the `KERNEL32 - VirtualAlloc` for windows systems.

## Is it ready

X **No.**

And Golang doesn't let me see into the heap, so I can't test some things properly,
like I could with C.

## Will it ever be?

X **I don't know if it's even worth it.**

Golang community even rejected https://github.com/golang/go/issues/51317 (proposal: arena: new package providing memory arenas #51317).
It seems like the community is very rejective on utilitarian approaches and is very conservative when it comes to safety.

Which I can understand to a certain point, but it's certainly just unproductive when you want to make something more performant that could benefit from having such arenas or ad-hoc memory allocation instead of having an automated "robot" we call "Garbage collector" do that for us. I as a programmer would definitely prefer to be given the controller sometimes and do whatever I want because maybe, I know better than the program sometimes.

So, the whole thing over "safety" made me, as the contributor of this repository, not so eager to complete this project.

The delve debugger simply doesn't support low-level heap debugging. 
Golang is trying actively to make your job harder when you're dealing with low level "unsafe code".
It's all just discouraging.

Maybe I will take a look once in a while on this repository just to take a look or to actually develop something.

But I promise nothing.

## TODO:

- [x] Support for allocation strategies & policies
- [x] Support for chunked memory allocation
- [x] Allocator doesn't allocate each time directly from CPU, but uses page-based chunks - and allocates new chunk per need.
- [ ] Support for arenas
- [ ] Support for resizing arenas (growing and shrinking)
- [ ] Removing unused arenas
- [ ] Support for `StringBuilder` type

## Supports

- Unix & POSIX compliant systems (Linux, macOS, ...)
- Windows

## Installation

```bash
go get -u github.com/exapsy/goumem
```

## Usage

### As simple as it goes

```go
package main

import (
    "fmt"

    "github.com/exapsy/goumem"
)

func main() {
    // Allocate 100 bytes
    block, err := goumem.Alloc(100)
    if err != nil {
        panic(err)
    }
	
    // Write to the memory
	block.Set("Hello World!")
	
    // Read from the memory
    var data string
    block.Get(&data)

    // Free the memory
    goumem.Free(mem)

    // Allocate 100 bytes
    mem, err = goumem.Alloc(100)
    if err != nil {
        panic(err)
    }

    // Free the memory
    goumem.Free(mem)
}
```

## Where we've:

### Seen vast improvements

- GC needs a lot of processing power in order to operate. Tracking objects, doing a lot of operations, etc. 
- GC needs a lot of memory in order to operate. It needs to keep track of all these objects, and it needs to keep them in memory until it's time to de-allocate them.
- GC while it's good at batch-deallocation when you don't want the "extra juice" and you don't care, when you actually care about the extra juice, it's not good at all. It will keep track of these objects until the end of time, and it will de-allocate them at the end of time. Which is not good if you want to de-allocate them at will.
- In order to do anything in GC there are a LOT of operations that need to be done. With a direct memory allocation, you just allocate and de-allocate. With GC, you need to allocate, keep track of the object, and then de-allocate. And that's just the tip of the iceberg. There are many more operations that need to be done in order to keep track of these objects, and that's why GC is slower by its nature.
- Pools are an especially good example of why you would need a direct memory allocation. You don't want to keep track of all these objects, you just want to allocate and de-allocate at will. And GC is not good at that. It will keep track of these objects until the end of time, and it will de-allocate them at the end of time. Which is not good if you want to de-allocate them at will.
- When you actually don't want all that tracking = all that extra CPU processing = all that extra memory
- You want to allocate and de-allocate at will
- Performance improvements and MUCH lower operations per nanosecond have been observed when using the non-GC memory allocation method.
- It's much more direct, and it's much more simple. It's just a memory allocation and de-allocation. No tracking, no nothing. 
- You don't have a powerful machine, and you want to squeeze the most out of it.

In these situations, **GC** is likely not the way to go.

**GC** is an autonomous system by its nature.

It that does things the way it wants to do them, and it's good at that, but it comes at a cost.

### Seen that GC is more optimal

Just to come clear.

**GC** is not bad, it's not the evil, it's actually very good at its job.

It's nice for its purpose, and it does it more than well.

That is also not necessarily bad or good. It's just a different way of doing things. And it's good to have both options.

For example a manual memory allocation is bad when

- You don't know what you're doing
- You don't know when/how/why to use it
- Memory leaks are probable to happen. Especially if you don't follow good practices.
- You don't want usually to batch de-allocate, for which GC is good at.
- Keeping track of objects, for which GC is good at.
- Thinking about memory allocation and de-allocation, for which GC is good at.
- Your program is short-lived and doesn't need much memory allocation and de-allocation, for which GC is good at. Like I would understand why a program like `jq` which decodes many amount of objects and probably does many allocations and de-allocations, would need it, but not why a program like `ls` would need it which usually prints a very small amount of objects and most likely doesn't do much memory allocation and de-allocation. Maybe bad example ... but you get my point. Don't rage at me at the issues about it.

## Benchmark

**Methodology**

Summing up **100000** `float64` matrices of **100x100**.

**Parameters**

```go
const (
    numMatrices = 100000
    rows        = 100
    cols        = 100
)
```

**Results**

```bash
$ go test -bench=. -benchmem -benchtime=240s
goos: linux
goarch: amd64
pkg: github.com/exapsy/goumem
cpu: 12th Gen Intel(R) Core(TM) i7-1255U
BenchmarkCustomMemory-12              10        25098256116 ns/op         400600 B/op      10001 allocs/op
BenchmarkGCMemory-12                   9        28630349082 ns/op       9229073352 B/op 10100069 allocs/op
PASS
ok      github.com/exapsy/goumem        561.607s
```
