package main

import "fmt"

func printSliceInfo(s []int) {
	fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}

func main() {
	var s []int
	for i := 0; i < 100; i++ {
		printSliceInfo(s)
		s = append(s, i)
	}

	s1 := []int{1, 2, 3, 4, 5, 6, 7}
	s2 := make([]int, 10)
	s3 := make([]int, 10, 32)
	printSliceInfo(s1)
	printSliceInfo(s2)
	printSliceInfo(s3)

	fmt.Println("copy slice")
	copy(s2, s1)
	fmt.Println(s1, s2)

	fmt.Println("delete elements from slice")
	s4 := append(s2[:3], s2[4:]...)
	fmt.Println(s4)
	printSliceInfo(s2)
	printSliceInfo(s4)

	fmt.Println("popping element from head")
	head := s2[0]
	s2 = s2[1:]
	printSliceInfo(s2)
	fmt.Println(head, s2)

	fmt.Println("popping from tail")
	tail := s2[len(s2)-1]
	s2 = s2[:len(s2)-1]
	printSliceInfo(s2)
	fmt.Println(tail, s2)

}
