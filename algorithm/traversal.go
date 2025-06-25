package algorithm

import "fmt"

func (node *Node) Traverse() {
	//if node == nil {
	//	return
	//}
	//fmt.Printf("node: %d\n", node.Val)
	//node.Left.Traverse()
	//node.Right.Traverse()
	node.TraverseFunc(func(node *Node) {
		node.Print()
	})
	fmt.Println()
}

func (node *Node) TraverseFunc(f func(*Node)) {
	if node == nil {
		return
	}
	node.Left.TraverseFunc(f)
	f(node)
	node.Right.TraverseFunc(f)
}

func (node *Node) TraverseWithChannel() chan *Node {
	out := make(chan *Node)
	go func() {
		node.TraverseFunc(func(node *Node) {
			out <- node
		})
		close(out)
	}()
	return out
}
