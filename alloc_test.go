package goumem

import (
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
	myArr, err := Alloc(10, MyStruct{})
	if err != nil {
		s.FailNow("Failed to allocate block")
	}

	myArr2, err := Alloc(20, MyStruct{})
	if err != nil {
		s.FailNow("Failed to allocate block")
	}

	s.Equal(uintptr(10*unsafe.Sizeof(MyStruct{})), myArr.Size())

	err = Free(myArr)
	if err != nil {
		s.FailNow("Failed to free allocated block")
	}

	err = Free(myArr2)
	if err != nil {
		s.FailNow("Failed to free allocated block")
	}
}

func TestAlloc(t *testing.T) {
	suite.Run(t, new(TestAllocSuite))
}
