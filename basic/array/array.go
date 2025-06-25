package main

import (
	"fmt"
)

func printArray(arr [5]int) {
	for i, v := range arr {
		fmt.Println(i, v)
	}
}

func printArrayPtr(arr *[5]int) {
	arr[0] = 100
	for i, v := range arr {
		fmt.Println(i, v)
	}
}

func split(start, end, interval int) (times []int) {
	if start > end || interval <= 0 {
		return
	}

	t := start
	for t < end {
		times = append(times, t)
		t += interval
	}
	times = append(times, end)
	return
}

func main() {
	var arr1 [5]int
	arr2 := [4]int{1, 2, 4, 5}
	arr3 := [...]int{1, 2, 13, 4, 5}
	grid := [...][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	fmt.Println(arr1, arr2, arr3)
	fmt.Println(grid)

	for i := 0; i < len(arr3); i++ {
		fmt.Println(arr3[i])
	}
	for i, v := range arr3 {
		fmt.Println(i, v)
	}

	maxi := -1
	maxValue := -1
	for i, v := range arr3 {
		if v > maxValue {
			maxi, maxValue = i, v
		}
	}
	fmt.Println(maxi, maxValue)

	sum := 0
	for _, v := range arr3 {
		sum += v
	}
	fmt.Println(sum)

	printArray(arr1)
	printArray(arr3)

	printArrayPtr(&arr1)
	printArray(arr1)

	fmt.Println(split(13, 13, 2))

	f := func() []int {
		return nil
	}
	fmt.Println(len(f()))
}
