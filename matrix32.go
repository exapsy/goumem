package goumem

import "unsafe"

type PointerMatrixFloat32 struct {
	virtualAddr uintptr
	rows        int
	cols        int
}

func NewMatrix32(rows, cols int) (*PointerMatrixFloat32, error) {
	virtualAddr, err := mem.Alloc(uintptr(rows * cols * 8))
	if err != nil {
		return nil, err
	}

	return &PointerMatrixFloat32{
		virtualAddr: virtualAddr,
		rows:        rows,
		cols:        cols,
	}, nil
}

func (ptr *PointerMatrixFloat32) Address() uintptr {
	return ptr.virtualAddr
}

func (ptr *PointerMatrixFloat32) Rows() int {
	return ptr.rows
}

func (ptr *PointerMatrixFloat32) Cols() int {
	return ptr.cols
}

func (ptr *PointerMatrixFloat32) Value() [][]float32 {
	matrix := make([][]float32, ptr.rows)
	for i := 0; i < ptr.rows; i++ {
		matrix[i] = make([]float32, ptr.cols)
		for j := 0; j < ptr.cols; j++ {
			matrix[i][j] = *(*float32)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*4))
		}
	}

	return matrix
}

func (ptr *PointerMatrixFloat32) Set(matrix [][]float32) {
	for i := 0; i < ptr.rows; i++ {
		for j := 0; j < ptr.cols; j++ {
			*(*float32)(unsafe.Pointer(ptr.virtualAddr + uintptr(i*ptr.cols+j)*4)) = matrix[i][j]
		}
	}
}

func (ptr *PointerMatrixFloat32) Free() error {
	return mem.Free(ptr.virtualAddr, uintptr(ptr.rows*ptr.cols*4))
}
