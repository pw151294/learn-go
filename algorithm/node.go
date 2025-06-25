package algorithm

import "fmt"

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

// 在C++里不能返回局部变量和临时变量的引用 因为局部/临时变量所在的空间会随着函数调用和代码块的离开所销毁
// Go语言不需要管局部变量的创建在堆上还是在栈上
func createNode(val int) *Node {
	return &Node{val, nil, nil}
}

// Print node *TreeNode声明的是函数的接收者
func (node *Node) Print() {
	fmt.Println(node.Val)
}

// TreeNode声明是值传递 会将node对象拷贝一份到函数体内 在该函数体内对于对象属性的操作不会被保存下来
func (node Node) setNewVal(newVal int) {
	node.Val = newVal
}

// SetVal *TreeNode声明是指针调用 是引用传递 函数体内对于对象的属性设置都会保存下来
func (node *Node) SetVal(val int) {
	if node == nil {
		fmt.Println("setting value to nil node")
	}
	node.Val = val
}
