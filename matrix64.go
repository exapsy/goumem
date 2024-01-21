package goumem

import (
	"fmt"
	"github.com/exapsy/goumem/allocator"
	"unsafe"
)

var (
	ErrMatrixZeroSize = fmt.Errorf("goumem: matrix size cannot be zero")
)

type PointerMatrixFloat64 struct {
	allocatedBlock *allocator.AllocatedBlock
	rows           int
	cols           int
}

func NewMatrixFloat64(rows, cols int) (*PointerMatrixFloat64, error) {
	block, err := mem.Alloc(uintptr((rows * cols) << 3))
	if err != nil {
		return nil, err
	}

	if rows == 0 || cols == 0 {
		return nil, ErrMatrixZeroSize
	}

	return &PointerMatrixFloat64{
		allocatedBlock: block,
		rows:           rows,
		cols:           cols,
	}, nil
}

func (ptr *PointerMatrixFloat64) Address() uintptr {
	return ptr.allocatedBlock.Addr()
}

func (ptr *PointerMatrixFloat64) Rows() int {
	return ptr.rows
}

func (ptr *PointerMatrixFloat64) Cols() int {
	return ptr.cols
}

func (ptr *PointerMatrixFloat64) Value() [][]float64 {
	matrix := make([][]float64, ptr.rows)
	flat := (*[1 << 30]float64)(unsafe.Pointer(ptr.allocatedBlock.Addr()))

	for i := range matrix {
		matrix[i] = flat[i*ptr.cols : (i+1)*ptr.cols]
	}

	return matrix
}

func (ptr *PointerMatrixFloat64) Set(matrix [][]float64) {
	if len(matrix) != ptr.rows {
		panic("goumem: matrix rows mismatch")
	}

	if len(matrix[0]) != ptr.cols {
		panic("goumem: matrix cols mismatch")
	}

	for i := 0; i < ptr.rows; i++ {
		for j := 0; j < ptr.cols; j++ {
			*(*float64)(unsafe.Pointer(ptr.allocatedBlock.Addr() + uintptr(i*ptr.cols+j)<<3)) = matrix[i][j]
		}
	}
}

func (ptr *PointerMatrixFloat64) Free() error {
	return mem.Free(ptr.allocatedBlock)
}

func (ptr *PointerMatrixFloat64) String() string {
	s := "["

	for i := 0; i < ptr.rows; i++ {
		for j := 0; j < ptr.cols; j++ {
			s += fmt.Sprintf("%f ", *(*float64)(unsafe.Pointer(ptr.allocatedBlock.Addr() + uintptr(i*ptr.cols+j)<<3)))
		}
		s += "\n"
	}

	s += "]"

	return s
}

type PoolMatrixFloat64 struct {
	totalMatrices int
	rows          int
	cols          int
	addr          uintptr
	matrices      []*PointerMatrixFloat64 // Preallocated matrix instances
	freeList      stack
}

type stack struct {
	matrices []*PointerMatrixFloat64
	len      int
}

func (s *stack) push(matrix *PointerMatrixFloat64) {
	if s.len >= len(s.matrices) {
		s.matrices = append(s.matrices, matrix)
	} else {
		s.matrices[s.len] = matrix
	}
	s.len++
}

func (s *stack) pop() *PointerMatrixFloat64 {
	if s.len == 0 {
		panic("goumem: stack is empty")
	}
	s.len--
	return s.matrices[s.len]
}

func NewPoolMatrix64(totalMatrices int, rows, cols int) (*PoolMatrixFloat64, error) {
	if rows == 0 || cols == 0 {
		return nil, ErrMatrixZeroSize
	}

	// Calculate the size of a matrix and the total size needed for all matrices
	matrixSize := uintptr(rows*cols) << 3
	totalSize := matrixSize * uintptr(totalMatrices)

	block, err := mem.Alloc(totalSize)
	if err != nil {
		return nil, err
	}

	matrices := make([]*PointerMatrixFloat64, totalMatrices)
	freeList := stack{matrices: make([]*PointerMatrixFloat64, totalMatrices)}
	for i := range matrices {
		matrices[i] = &PointerMatrixFloat64{
			allocatedBlock: block,
			rows:           rows,
			cols:           cols,
		}
		freeList.push(matrices[i])
	}

	return &PoolMatrixFloat64{
		totalMatrices: totalMatrices,
		rows:          rows,
		cols:          cols,
		addr:          block.Addr(),
		matrices:      matrices,
		freeList:      freeList,
	}, nil
}

func (pool *PoolMatrixFloat64) Get() *PointerMatrixFloat64 {
	return pool.freeList.pop()
}

func (pool *PoolMatrixFloat64) Free(matrix *PointerMatrixFloat64) {
	pool.freeList.push(matrix)
}
