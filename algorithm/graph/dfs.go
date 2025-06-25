package main

import (
	"fmt"
)

type Node struct {
	Value    int
	Children []*Node
}

func DFS(node *Node) []int {
	if node == nil {
		return nil
	}
	result := make([]int, 0)
	result = append(result, node.Value)
	for _, child := range node.Children {
		childResult := DFS(child)
		result = append(result, childResult...)
	}
	return result
}

func main() {
	// 构造示例树
	root := &Node{
		Value: 1,
		Children: []*Node{
			&Node{Value: 2, Children: []*Node{&Node{Value: 4}, &Node{Value: 5}}},
			&Node{Value: 3, Children: []*Node{&Node{Value: 6}, &Node{Value: 7}}},
		},
	}

	// 调用DFS并打印结果
	result := DFS(root)
	fmt.Println(result) // 输出: [1 2 4 5 3 6 7]
}
