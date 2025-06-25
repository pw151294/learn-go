package main

import "fmt"

func updateSlice(s []int) {
	s[0] = 100
}

func printSlice(arr []int) {
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
}

func add(e int, arr *[]int) {
	*arr = append(*arr, e)
}

func modify(idx, val int, arr *[]int) {
	list := *arr
	list[idx] = val
}

// 求集合的并集
func union(arr1, arr2 []int) []int {
	set := make(map[int]struct{})
	for _, val := range arr1 {
		set[val] = struct{}{}
	}

	result := arr1
	for _, val := range arr2 {
		if _, ok := set[val]; !ok {
			set[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

// 求集合的交集
func intersection(arr1, arr2 []int) (result []int) {
	set := make(map[int]struct{})
	for _, val := range arr1 {
		set[val] = struct{}{}
	}

	for _, val := range arr2 {
		if _, ok := set[val]; ok {
			result = append(result, val)
		}
	}

	return
}

func copySlice(s []int) []int {
	cs := make([]int, 0, len(s))
	for _, v := range s {
		cs = append(cs, v)
	}
	return cs
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("arr[2,6]=", arr[2:6])
	fmt.Println("arr[:6]=", arr[:6])
	fmt.Println("arr[:]=", arr[:])
	fmt.Println("arr[2:]=", arr[2:])
	add(8, &arr)
	modify(7, 10, &arr)
	fmt.Println("arr=", arr)

	s1 := arr[2:]
	s2 := arr[:]
	fmt.Printf("begin to update s1 and s2")
	updateSlice(s1)
	updateSlice(s2)
	fmt.Println(s1, s2, arr)
	printSlice(arr)

	s3 := arr[:]
	fmt.Println("begin to re slice")
	s3 = s3[:5]
	s3 = s3[2:]
	fmt.Println(s3)

	fmt.Println("extending slice")
	s1 = arr[2:6]
	s2 = s1[3:5]
	s3 = append(s1, 10)
	s4 := append(s3, 11)
	s5 := append(s4, 12)
	fmt.Printf("s1=%v len(s1)=%d cap(s1)=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2=%v len(s2)=%d cap(s2)=%d\n", s2, len(s2), cap(s2))
	fmt.Println(s1, s2, arr)
	fmt.Println(s3, s4, s5)

	fmt.Println(union(s3, union(s4, s5)))
	fmt.Print(intersection(s3, intersection(s4, s5)))
	fmt.Println("copied slice", copySlice([]int{1, 2, 3, 4, 5}))
}
