package goumem

import (
	"math/rand"
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
	pool, err := NewPoolMatrix64(numMatrices, rows, cols)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var r float64
		for pb.Next() {
			randVal := rand.Float64()
			for i := 0; i < numMatrices; i++ {
				matrix := pool.Get()
				r += simulateReadWrite(matrix, randVal)
				pool.Free(matrix)
			}
		}
		result = r
	})
	b.StopTimer()
	//poolMatrix64, err := NewPoolMatrix64(numMatrices, rows, cols)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//
	//var r float64
	//for i := 0; i < b.N; i++ {
	//	randVal := rand.Float64()
	//	for j := 0; j < numMatrices; j++ {
	//		matrix := poolMatrix64.Get()
	//		r += simulateReadWrite(matrix, randVal)
	//		poolMatrix64.Free(matrix)
	//	}
	//	runtime.GC()
	//}
	//b.StopTimer()          // Stop the timer before doing operations not related to the actual benchmark
	//result = r             // Assign the final result to the global variable
	//_ = fmt.Sprint(result) // Print the result or use it in a way that ensures it's not optimized out
}

func BenchmarkGCMemory(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var r float64
		for pb.Next() {
			matrix := makeSampleMatrix(rows, cols)
			randVal := rand.Float64()
			for i := 0; i < numMatrices; i++ {
				r += simulateReadWriteGC(matrix, randVal)
			}
		}
		result = r
	})
	b.StopTimer()
	//var r float64
	//b.ResetTimer()
	//matrices := make([][]float64, numMatrices)
	//for i := 0; i < b.N; i++ {
	//	randVal := rand.Float64()
	//	for j := 0; j < numMatrices; j++ {
	//		matrix := makeSampleMatrix(rows, cols)
	//		r += simulateReadWriteGC(matrix, randVal)
	//		matrices = matrix
	//	}
	//	runtime.GC()
	//}
	//runtime.KeepAlive(matrices)
	//b.StopTimer()          // Stop the timer before doing operations not related to the actual benchmark
	//result = r             // Assign the final result to the global variable
	//_ = fmt.Sprint(result) // Print the result or use it in a way that ensures it's not optimized out
}

func simulateReadWrite(matrix *PointerMatrixFloat64, randVal float64) float64 {
	var sum float64
	var val float64
	rows := matrix.Rows()
	cols := matrix.Cols()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val = *(*float64)(unsafe.Pointer(matrix.Address() + uintptr(i*cols+j)*8))
			val *= randVal
			*(*float64)(unsafe.Pointer(matrix.Address() + uintptr(i*cols+j)*8)) = val
			sum += val
		}
	}
	return sum
}

func simulateReadWriteGC(matrix [][]float64, randVal float64) float64 {
	var sum float64
	for r := range matrix {
		for c := range matrix[r] {
			val := matrix[r][c]
			val *= randVal
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
