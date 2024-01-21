package goumem

import (
	"github.com/exapsy/goumem/allocator"
	"reflect"
)

func Alloc(t interface{}) (*allocator.AllocatedBlock, error) {
	tt := reflect.TypeOf(t)
	b, err := mem.Alloc(tt.Size())
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Free(block *allocator.AllocatedBlock) error {
	return mem.Free(block)
}

func init() {
	if mem == nil {
		mem = allocator.Default()
	}
}
