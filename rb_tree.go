package main

import (
	"fmt"
)

// Properties of red-black tree:
// 1. Every node is either red or black.
// 2. The root of the tree is always black.
// 3. Red nodes cannot have red children.
// 4. Every path from a node to its descendant null nodes
//    (leaves) has the same number of black nodes.
// 5. All leaves (nil nodes) are black.

type Color bool

const (
	Red   Color = true
	Black       = false
)

type Node struct {
	color               Color
	data                int
	left, right, parent *Node
}

type RBTree struct {
	root *Node
}

func NewNode(data int) *Node {
	return &Node{data: data, color: Red, left: nil, right: nil, parent: nil}
}

func NewRedBlackTree() *RBTree {
	return &RBTree{root: nil}
}

// TODO: func (rbt *RBTree) Delete(data int) bool   { return false }
func (rbt *RBTree) LeftRotate(n *Node) {
	y := n.right
	n.right = y.left

	if y.left != nil {
		y.left.parent = n
	}
	y.parent = n.parent
	if n.parent == nil {
		rbt.root = y
	} else if n == n.parent.left {
		n.parent.left = y
	} else {
		n.parent.right = y
	}
	y.left = n
	n.parent = y
}
func (rbt *RBTree) RightRotate(n *Node) {
	y := n.left
	n.left = y.right

	if y.right != nil {
		y.right.parent = n
	}
	y.parent = n.parent
	if n.parent == nil {
		rbt.root = y
	} else if n == n.parent.right {
		n.parent.right = y
	} else {
		n.parent.left = y
	}
	y.right = n
	n.parent = y
}
func (rbt *RBTree) FixInsert(n *Node) {
	for n != rbt.root && n.parent.color == Red {
		parent := n.parent
		grandparent := parent.parent

		// if parent == nil || grandparent == nil {
		// 	fmt.Println("ultra giga hughmongous cringe")
		// 	break
		// }

		if parent == grandparent.left {
			uncle := grandparent.right
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle.color == Red {
				parent.color = Black
				uncle.color = Black
				grandparent.color = Red
				n = grandparent
			} else {
				// Case 2.1 (Uncle is black and n is the right child):
				//   Perform a left rotation on parent
				if n == parent.right {
					n = parent
					rbt.LeftRotate(n)
				}
				// Case 2.2 (Uncle is black and n is the left child):
				//   Recolor parent and grandparent, right rotate on grandparent
				parent.color = Black
				grandparent.color = Red
				rbt.RightRotate(grandparent)
			}
		} else {
			uncle := grandparent.left
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle.color == Red {
				parent.color = Black
				uncle.color = Black
				grandparent.color = Red
				n = grandparent
			} else {
				// Case 2 (Uncle is black):
				if n == parent.left {
					n = parent
					rbt.RightRotate(n)
				}
				parent.color = Black
				grandparent.color = Red
				rbt.LeftRotate(grandparent)
			}
		}
	}
	rbt.root.color = Black
}

func (rbt *RBTree) Insert(data int) bool {
	if rbt.root == nil {
		rbt.root = &Node{data: data, color: Black}
		return true
	}
	new_node := NewNode(data)
	var parent *Node
	current := rbt.root

	// BTS insert
	for current != nil {
		parent = current
		if new_node.data < current.data {
			current = current.left
		} else if new_node.data > current.data {
			current = current.right
		} else {
			// Duplicate
			return false
		}
	}
	new_node.parent = parent
	if parent.data < new_node.data {
		parent.right = new_node
	} else {
		parent.left = new_node
	}

	rbt.FixInsert(new_node)
	return true
}

// TODO: func (rbt *RBTree) Search(data int) bool {}
func Preorder(n *Node) {
	if n != nil {
		fmt.Print(n.data, " ")
		Preorder(n.left)
		Preorder(n.right)
	}
}
func Postorder(n *Node) {
	if n != nil {
		Preorder(n.left)
		fmt.Print(n.data, " ")
		Preorder(n.right)
	}
}

func Inorder(n *Node) {
	if n != nil {
		Preorder(n.left)
		Preorder(n.right)
		fmt.Print(n.data, " ")
	}
}

func main() {
	rbt := NewRedBlackTree()
	rbt.Insert(10)
	rbt.Insert(11)
	rbt.Insert(9)
	rbt.Insert(9)
	rbt.Insert(15)
	rbt.Insert(35)
	rbt.Insert(25)
	rbt.Insert(20)
	rbt.Insert(1)
	Inorder(rbt.root)
	fmt.Println("hello")
}
