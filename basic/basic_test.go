package main

import (
	"math"
	"testing"
)

func CalcTriangle(a, b int) int {
	return int(math.Sqrt(float64(a*a + b*b)))
}

// 表格驱动测试
func TestTriangle(t *testing.T) {
	tests := []struct{ x, y, z int }{
		{3, 4, 5},
		{5, 12, 13},
		{8, 15, 17},
		{30000, 40000, 50000},
	}

	for _, tt := range tests {
		if actual := CalcTriangle(tt.x, tt.y); actual != tt.z {
			t.Errorf("CalcTriangle(%d, %d) = %d, want %d", tt.x, tt.y, tt.z, actual)
		}
	}
}
