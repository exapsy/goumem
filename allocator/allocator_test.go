package allocator

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AllocatorTestSuite struct {
	suite.Suite
	allocator MemoryAllocator
}

func (suite *AllocatorTestSuite) SetupTest() {
	suite.allocator = Default()
}

func (suite *AllocatorTestSuite) TestSetGet() {
	data := "test data"
	block, err := suite.allocator.Alloc(uintptr(len(data)))
	if err != nil {
		suite.FailNow("Failed to allocate block", err)
	}

	suite.Equal(block.IsFreed(), false)

	block.Set(data)

	var got string
	block.Get(&got)

	suite.Equal(data, got)

	err = suite.allocator.Free(block)
	if err != nil {
		suite.FailNow("Failed to free allocated block", err)
	}

	suite.Equal(block.IsFreed(), true)
}

func TestAllocatorTestSuite(t *testing.T) {
	suite.Run(t, new(AllocatorTestSuite))
}
