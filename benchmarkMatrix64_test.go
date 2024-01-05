package goumem

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"unsafe"
)

// Global variable to prevent compiler optimizations
var result float64

// In total 100000 matrices of 100x100 elements
// Each element is a float64
// Which means that each matrix is 80000 bytes
// In total 8000000000 bytes
// Which means 7.450580596923828 GB
const (
	numMatrices = 100000
	rows        = 100
	cols        = 100
)

func BenchmarkCustomMemory(b *testing.B) {
	b.ResetTimer()
	pool, err := NewPoolMatrix64(numMatrices, rows, cols)
	if err != nil {
		b.Fatal(err)
	}

	var r float64
	for i := 0; i < b.N; i++ {
		for j := 0; j < numMatrices; j++ {
			matrix := pool.GetMatrix()
			r += simulateReadWrite(matrix)
			pool.Free(matrix)
			runtime.KeepAlive(matrix)
		}
	}
	b.StopTimer()          // Stop the timer before doing operations not related to the actual benchmark
	result = r             // Assign the final result to the global variable
	_ = fmt.Sprint(result) // Print the result or use it in a way that ensures it's not optimized out
}

func BenchmarkGCMemory(b *testing.B) {
	var r float64
	b.ResetTimer()
	matrices := make([][]float64, numMatrices)
	for i := 0; i < b.N; i++ {
		for j := 0; j < numMatrices; j++ {
			matrix := makeSampleMatrix(rows, cols)
			r += simulateReadWriteGC(matrix)
			matrices = matrix
		}
	}
	runtime.KeepAlive(matrices)
	b.StopTimer()          // Stop the timer before doing operations not related to the actual benchmark
	result = r             // Assign the final result to the global variable
	_ = fmt.Sprint(result) // Print the result or use it in a way that ensures it's not optimized out
}

func simulateReadWrite(matrix *PointerMatrixFloat64) float64 {
	var sum float64
	var val float64
	rows := matrix.Rows()
	cols := matrix.Cols()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val = *(*float64)(unsafe.Pointer(matrix.Address() + uintptr(i*cols+j)*8))
			val *= rand.Float64()
			*(*float64)(unsafe.Pointer(matrix.Address() + uintptr(i*cols+j)*8)) = val
			sum += val
		}
	}
	return sum
}

func simulateReadWriteGC(matrix [][]float64) float64 {
	var sum float64
	for r := range matrix {
		for c := range matrix[r] {
			val := matrix[r][c]
			val *= rand.Float64()
			matrix[r][c] = val
			sum += val
		}
	}
	return sum
}

func makeSampleMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()
		}
	}
	return matrix
}
