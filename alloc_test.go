package goumem

import (
	"github.com/exapsy/goumem/allocator"
	"github.com/stretchr/testify/suite"
	"testing"
	"unsafe"
)

type TestAllocSuite struct {
	suite.Suite
}

func (s *TestAllocSuite) SetupSuite() {
}

func (s *TestAllocSuite) TestAlloc() {
	type MyStruct struct {
		a int
		b float64
		c string
	}
	myArr, err := Alloc([10]MyStruct{})
	if err != nil {
		s.FailNow("Failed to allocate block")
	}

	s.Equal(uintptr(unsafe.Sizeof([10]MyStruct{})), myArr.Size())

	myArr2, err := Alloc([20]MyStruct{})
	if err != nil {
		s.FailNow("Failed to allocate block")
	}

	s.Equal(uintptr(unsafe.Sizeof([20]MyStruct{})), myArr2.Size())

	*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())) = [10]MyStruct{{8, 2.0, ""}}
	s.Equal(8, (*(*[10]MyStruct)(unsafe.Pointer(myArr.Addr())))[0].a)

	err = Free(myArr)
	if err != nil {
		s.FailNow("Failed to free allocated block")
	}

	err = Free(myArr)
	s.EqualError(err, allocator.ErrAllocatedBlockAlreadyFreed.Error())
	s.Equal(myArr.IsFreed(), true)

	err = Free(myArr2)
	if err != nil {
		s.FailNow("Failed to free allocated block")
	}
}

func TestAlloc(t *testing.T) {
	suite.Run(t, new(TestAllocSuite))
}
