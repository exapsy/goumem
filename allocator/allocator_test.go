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
	// Add any setup code here...
}

func (suite *AllocatorTestSuite) TestSetGet() {
	data := "test data"
	block, err := suite.allocator.Alloc(uintptr(len(data)))
	if err != nil {
		suite.FailNow("Failed to allocate block", err)
	}

	block.Set(data)

	var got string
	block.Get(&got)

	suite.Equal(data, got)
}

func TestAllocatorTestSuite(t *testing.T) {
	suite.Run(t, new(AllocatorTestSuite))
}
