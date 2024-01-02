package goumem

import "unsafe"

type PointerMatrixFloat64 struct {
	virtualAddr uintptr
	rows        int
	cols        int
}

func NewMatrixFloat64(rows, cols int) (*PointerMatrixFloat64, error) {
	virtualAddr, err := mem.Alloc(uintptr(rows * cols * 8))
	if err != nil {
		return nil, err
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
	for i := 0; i < ptr.rows; i++ {
		for j := 0; j < ptr.cols; j++ {
			*(*float64)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*8)) = matrix[i][j]
		}
	}
}

func (ptr *PointerMatrixFloat64) Free() error {
	return mem.Free(ptr.virtualAddr, uintptr(ptr.rows*ptr.cols*8))
}
