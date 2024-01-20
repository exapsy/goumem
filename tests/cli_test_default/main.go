package main

import (
	"fmt"
	"github.com/exapsy/goumem/allocator"
)

var (
	mem             allocator.MemoryAllocator
	allocatedBlocks int
	arrayOfBlocks   []allocator.AllocatedBlock
)

func main() {
	fmt.Println(`Commands: 

alloc <size> - allocates <number> bytes and prints the index of the allocated block
free <number> - frees the block with the given index
exit - exits the program
`)

	for {
		command, err := scanCommand()
		if err != nil {
			fmt.Println("invalid command: ", err)
			continue
		}

		switch command {
		case "alloc":
			handleAlloc()

		case "free":
			handleFree()
		}
	}
}

func scanCommand() (string, error) {
	var command string

	_, err := fmt.Scan(&command)
	if err != nil {
		return "", fmt.Errorf("invalid command: %w", err)
	}

	return command, nil
}

func handleAlloc() {
	var size uintptr
	_, err := fmt.Scanln(&size)
	if err != nil {
		fmt.Println("invalid size: ", err)
		return
	}

	block, err := mem.Alloc(size)
	if err != nil {
		fmt.Println("allocation error: ", err)
		return
	}

	fmt.Println("index: ", allocatedBlocks)
	fmt.Printf("allocated block: %+v\n", block)
	fmt.Println()

	arrayOfBlocks = append(arrayOfBlocks, *block)
	allocatedBlocks++
}

func handleFree() {
	var blockIndex uintptr
	_, err := fmt.Scan(&blockIndex)
	if err != nil {
		fmt.Println("invalid block index: ", err)
		return
	}

	if blockIndex >= uintptr(allocatedBlocks) {
		fmt.Println("invalid block index: ", err)
		return
	}

	err = mem.Free(&arrayOfBlocks[blockIndex])

	fmt.Printf("block %d freed", blockIndex)
	fmt.Println("block: ", arrayOfBlocks[blockIndex])
	fmt.Println()

	arrayOfBlocks = append(arrayOfBlocks[:blockIndex], arrayOfBlocks[blockIndex+1:]...)
	allocatedBlocks--
}

func init() {
	mem = allocator.Default()
	allocatedBlocks = 0
	arrayOfBlocks = []allocator.AllocatedBlock{}
}
