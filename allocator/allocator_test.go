package allocator

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"unsafe"
)

type AllocatorTestSuite struct {
	suite.Suite
	allocator MemoryAllocator
}

func (suite *AllocatorTestSuite) SetupTest() {
	suite.allocator = Default()
}

func (suite *AllocatorTestSuite) TestAllocFree() {
	block, err := suite.allocator.Alloc(10)
	if err != nil {
		suite.FailNow("Failed to allocate block", err)
	}

	suite.Equal(uintptr(10), block.Size())

	err = suite.allocator.Free(block)
	if err != nil {
		suite.FailNow("Failed to free allocated block", err)
	}
}

func (suite *AllocatorTestSuite) TestSet() {
	suite.Run("struct", func() {
		type MyStruct struct {
			a int
			b float64
			c string
		}
		myArr, err := suite.allocator.Alloc(unsafe.Sizeof([10]MyStruct{}))
		if err != nil {
			suite.FailNow("Failed to allocate block", err)
		}

		Set[[10]MyStruct](myArr, [10]MyStruct{{8, 3.14, "test data"}})
		suite.Equal(8, (*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())))[0].a)
		suite.Equal(3.14, (*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())))[0].b)
		suite.Equal("test data", (*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())))[0].c)

		err = suite.allocator.Free(myArr)
		if err != nil {
			suite.FailNow("Failed to free allocated block", err)
		}

		err = suite.allocator.Free(myArr)
		suite.EqualError(err, ErrAllocatedBlockAlreadyFreed.Error())
		suite.Equal(myArr.IsFreed(), true)
	})
	suite.Run("string", func() {
		data := "test data"
		block, err := suite.allocator.Alloc(uintptr(len(data)))
		if err != nil {
			suite.FailNow("Failed to allocate block", err)
		}

		suite.Equal(block.IsFreed(), false)

		Set(block, data)

		var got string
		*(*string)(unsafe.Pointer(&got)) = *(*string)(unsafe.Pointer(block.Addr()))

		suite.Equal(data, got)

		err = suite.allocator.Free(block)
		if err != nil {
			suite.FailNow("Failed to free allocated block", err)
		}

		suite.Equal(block.IsFreed(), true)
	})
}

func (suite *AllocatorTestSuite) TestGet() {
	suite.Run("struct", func() {
		type MyStruct struct {
			a int
			b float64
			c string
		}
		myArr, err := suite.allocator.Alloc(unsafe.Sizeof([10]MyStruct{}))
		if err != nil {
			suite.FailNow("Failed to allocate block", err)
		}

		var got [10]MyStruct

		*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())) = [10]MyStruct{{14, 3.14, "bruh"}}
		got = *(*[10]MyStruct)(unsafe.Pointer(myArr.Addr()))
		suite.Equal(14, got[0].a)
		suite.Equal(3.14, got[0].b)
		suite.Equal("bruh", got[0].c)

		got = Get[[10]MyStruct](myArr)

		suite.Equal(14, got[0].a)
		suite.Equal(3.14, got[0].b)
		suite.Equal("bruh", got[0].c)

		err = suite.allocator.Free(myArr)
		if err != nil {
			suite.FailNow("Failed to free allocated block", err)
		}
	})
	suite.Run("str_array", func() {
		data := [3]string{"test", "data", "lol"}
		block, err := suite.allocator.Alloc(uintptr(unsafe.Sizeof(data)))
		if err != nil {
			suite.FailNow("Failed to allocate block", err)
		}

		*(*[3]string)(unsafe.Pointer(block.Addr())) = data
		testGot := *(*[3]string)(unsafe.Pointer(block.Addr()))
		suite.Equal(testGot, *(*[3]string)(unsafe.Pointer(block.Addr())))

		got := Get[[3]string](block)

		suite.Equal(data, got)

		err = suite.allocator.Free(block)
		if err != nil {
			suite.FailNow("Failed to free allocated block", err)
		}
	})
	suite.Run("string", func() {
		data := "test data"
		block, err := suite.allocator.Alloc(uintptr(len(data)))
		if err != nil {
			suite.FailNow("Failed to allocate block", err)
		}

		*(*string)(unsafe.Pointer(block.Addr())) = data
		testGot := *(*string)(unsafe.Pointer(block.Addr()))
		suite.Equal(testGot, *(*string)(unsafe.Pointer(block.Addr())))

		got := Get[string](block)

		suite.Equal(data, got)

		err = suite.allocator.Free(block)
		if err != nil {
			suite.FailNow("Failed to free allocated block", err)
		}
	})
}

func (suite *AllocatorTestSuite) TestCopy() {
	data := "test data"
	srcBlock, err := suite.allocator.Alloc(uintptr(len(data)))
	if err != nil {
		suite.FailNow("Failed to allocate block", err)
	}

	dstBloc, err := suite.allocator.Alloc(uintptr(len(data)))
	if err != nil {
		suite.FailNow("Failed to allocate block", err)
	}

	// Set data in srcBlock
	Set[string](srcBlock, data)

	// Copy data from srcBlock to dstBlock
	err = suite.allocator.Copy(dstBloc, srcBlock)

	// Get data from dstBlock
	got := Get[string](srcBlock)

	suite.Equal(data, got)

	err = suite.allocator.Free(srcBlock)
	if err != nil {
		suite.FailNow("Failed to free allocated block", err)
	}

	err = suite.allocator.Free(dstBloc)
	if err != nil {
		suite.FailNow("Failed to free allocated block", err)
	}
}

func TestAllocatorTestSuite(t *testing.T) {
	suite.Run(t, new(AllocatorTestSuite))
}

func ExampleMemoryAllocator_Alloc() {
	// Allocate a block of memory
	block, err := Default().Alloc(10)
	if err != nil {
		panic(err)
	}

	// Free the block
	err = Default().Free(block)
	if err != nil {
		panic(err)
	}
}

func ExampleAllocatedBlock_Set() {
	// Allocate a block of memory
	block, err := Default().Alloc(10)
	if err != nil {
		panic(err)
	}

	// Set data in the block
	Set[string](block, "test data")

	// Free the block
	err = Default().Free(block)
	if err != nil {
		panic(err)
	}
}

func ExampleAllocatedBlock_Get() {
	// Allocate a block of memory
	block, err := Default().Alloc(10)
	if err != nil {
		panic(err)
	}

	// Set data in the block
	Set(block, "test data")

	// Get data from the block
	got := Get[string](block)
	print(got)

	// Free the block
	err = Default().Free(block)
	if err != nil {
		panic(err)
	}
}

func ExampleAllocatedBlock_Copy() {
	// Allocate a block of memory
	srcBlock, err := Default().Alloc(10)
	if err != nil {
		panic(err)
	}

	// Allocate a block of memory
	dstBlock, err := Default().Alloc(10)
	if err != nil {
		panic(err)
	}

	// Set data in the srcBlock
	Set(srcBlock, "test data")

	// Copy data from srcBlock to dstBlock
	err = Default().Copy(dstBlock, srcBlock)
	if err != nil {
		panic(err)
	}

	// Free the srcBlock
	err = Default().Free(srcBlock)
	if err != nil {
		panic(err)
	}

	// Free the dstBlock
	err = Default().Free(dstBlock)
	if err != nil {
		panic(err)
	}
}
