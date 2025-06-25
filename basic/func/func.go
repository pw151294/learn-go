package main

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
)

func div1(a, b int) (int, int) {
	return a / b, a % b
}

// 函数体较长的场景下不建议使用这样的返回值定义方法
func div2(a, b int) (q, r int) {
	q = a / b
	r = a % b
	return
}

func div3(a, b int) (q, r int) {
	return a / b, a % b
}

func calc1(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		q, _ := div3(a, b) //不需要使用的变量 直接用_代替 否在被定义的变量不使用的话会报错
		return q
	default:
		panic("invalid operation")
	}
}

func calc2(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		q, _ := div3(a, b) //不需要使用的变量 直接用_代替 否在被定义的变量不使用的话会报错
		return q, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", op)
	}
}

func apply(op func(int, int) int, a, b int) int {
	ptr := reflect.ValueOf(op).Pointer()
	opName := runtime.FuncForPC(ptr).Name()
	fmt.Printf("Calling function %s with arguments (%d, %d)", opName, a, b)
	return op(a, b)
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func swapByValue(a, b int) (int, int) {
	return b, a
}

func swapByRef(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	fmt.Println(div1(13, 3))

	q, r := div2(13, 3)
	fmt.Println(q, r)

	q1, r1 := div3(13, 3)
	fmt.Println(q1, r1)

	fmt.Println(calc1(3, 4, "/"))
	fmt.Println(calc2(3, 4, "!"))

	fmt.Println(apply(pow, 3, 4))
	fmt.Println(apply(
		func(a int, b int) int {
			return a + b
		}, 3, 4))

	fmt.Printf("%d", sum(1, 2, 3))

	a, b := 3, 4
	a, b = swapByValue(a, b)
	fmt.Printf("swapByValue: %d %d", a, b)
	swapByRef(&a, &b)
	fmt.Printf("swapByRef: %d %d", a, b)
}
