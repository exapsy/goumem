package goumem

import (
	"fmt"
	"unsafe"
)

type PointerMatrixFloat64 struct {
	virtualAddr uintptr
	rows        int
	cols        int
}

func NewMatrix64(rows, cols int) (*PointerMatrixFloat64, error) {
	virtualAddr, err := mem.Alloc(uintptr(rows * cols * 8))
	if err != nil {
		return nil, err
	}

	if rows == 0 || cols == 0 {
		return nil, ErrMatrixZeroSize
	}

	return &PointerMatrixFloat64{
		virtualAddr: virtualAddr,
		rows:        rows,
		cols:        cols,
	}, nil
}

func (ptr *PointerMatrixFloat64) Address() uintptr {
	return ptr.virtualAddr
}

func (ptr *PointerMatrixFloat64) Rows() int {
	return ptr.rows
}

func (ptr *PointerMatrixFloat64) Cols() int {
	return ptr.cols
}

func (ptr *PointerMatrixFloat64) Value() [][]float64 {
	matrix := make([][]float64, ptr.rows)
	for i := 0; i < ptr.rows; i++ {
		matrix[i] = make([]float64, ptr.cols)
		for j := 0; j < ptr.cols; j++ {
			matrix[i][j] = *(*float64)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*8))
		}
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
			*(*float64)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*8)) = matrix[i][j]
		}
	}
}

func (ptr *PointerMatrixFloat64) Free() error {
	return mem.Free(ptr.virtualAddr, uintptr(ptr.rows*ptr.cols*8))
}

func (ptr *PointerMatrixFloat64) String() string {
	s := "["

	for i := 0; i < ptr.rows; i++ {
		for j := 0; j < ptr.cols; j++ {
			s += fmt.Sprintf("%f ", *(*float64)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*8)))
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
	matrixSize := uintptr(rows * cols * 8)
	totalSize := matrixSize * uintptr(totalMatrices)

	// Add extra space for alignment
	totalSize += 63 // 63 is 64-1, which is the maximum offset we might need to align the memory

	addr, err := mem.Alloc(totalSize)
	if err != nil {
		return nil, err
	}

	// Calculate the aligned address
	alignedAddr := (addr + 63) &^ uintptr(63)

	matrices := make([]*PointerMatrixFloat64, totalMatrices)
	freeList := stack{matrices: make([]*PointerMatrixFloat64, totalMatrices)}
	for i := range matrices {
		matrices[i] = &PointerMatrixFloat64{
			virtualAddr: alignedAddr + uintptr(i)*matrixSize,
			rows:        rows,
			cols:        cols,
		}
		freeList.push(matrices[i])
	}

	return &PoolMatrixFloat64{
		totalMatrices: totalMatrices,
		rows:          rows,
		cols:          cols,
		addr:          alignedAddr,
		matrices:      matrices,
		freeList:      freeList,
	}, nil
}

func (pool *PoolMatrixFloat64) GetMatrix() *PointerMatrixFloat64 {
	return pool.freeList.pop()
}

func (pool *PoolMatrixFloat64) Free(matrix *PointerMatrixFloat64) {
	pool.freeList.push(matrix)
}
