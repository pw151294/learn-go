package main

import (
	"fmt"
	"iflytek.com/weipan4/learn-go/algorithm"
)

type SysTreeNode struct {
	node *algorithm.Node
}

func (sysNode *SysTreeNode) postOrder() {
	if sysNode == nil || sysNode.node == nil {
		return
	}
	left := SysTreeNode{sysNode.node.Left}
	left.postOrder()
	right := SysTreeNode{sysNode.node.Right}
	right.postOrder()
	sysNode.node.Print()
}

type EmbeddingTreeNode struct {
	*algorithm.Node //使用内嵌的方式扩展已有的类型
}

func (embeddingNode *EmbeddingTreeNode) postOrder() {
	if embeddingNode == nil || embeddingNode.Node == nil {
		return
	}
	left := EmbeddingTreeNode{embeddingNode.Node.Left}
	left.postOrder()
	right := EmbeddingTreeNode{embeddingNode.Node.Right}
	right.postOrder()
	embeddingNode.Print()
}

func (embeddingNode *EmbeddingTreeNode) Traverse() {
	if embeddingNode == nil || embeddingNode.Node == nil {
		return
	}
	fmt.Println(embeddingNode.Node.Val)
	left := EmbeddingTreeNode{embeddingNode.Node.Left}
	left.Traverse()
	right := EmbeddingTreeNode{embeddingNode.Node.Right}
	right.Traverse()
}

func main() {
	var root algorithm.Node
	fmt.Println(root)

	root = algorithm.Node{Val: 1}
	root.Left = &algorithm.Node{}
	root.Right = &algorithm.Node{Val: 2}
	root.Right.Left = new(algorithm.Node)
	//root.Left.Right =
	fmt.Println(root)
	fmt.Println(root.Right)
	fmt.Println(root.Left)
	nodes := []algorithm.Node{
		{Val: 1},
		{Val: 2},
		{Val: 3, Left: nil, Right: &root},
	}
	fmt.Println(nodes)

	root.Print()
	root.SetVal(4)
	root.Print()
	root.Traverse()

	pRoot := &root
	pRoot.Print()
	pRoot.SetVal(40)
	pRoot.Print()

	sysRoot := SysTreeNode{pRoot}
	sysRoot.postOrder()

	embeddingRoot := EmbeddingTreeNode{pRoot}
	embeddingRoot.postOrder()
	embeddingRoot.Traverse()

	//通过函数式编程查询树的节点个数
	nodeCnt := 0
	root.TraverseFunc(func(node *algorithm.Node) {
		nodeCnt++
		node.Print()
	})
	fmt.Printf("node count: %d\n", nodeCnt)

	c := root.TraverseWithChannel()
	maxNode := 0
	for node := range c {
		maxNode = max(node.Val, maxNode)
	}
	fmt.Printf("max node: %d\n", maxNode)
}
