package main

import (
	"fmt"
	"math"
)

// Properties of red-black tree:
// 1. Every node is either red or black.
// 2. The root of the tree is always black.
// 3. Red nodes cannot have red children.
// 4. Every path from a node to its descendant null nodes
//    (leaves) has the same number of black nodes.
// 5. All leaves (nil nodes) are black.

const colorRed = "\033[0;31m"
const colorNone = "\033[0m"

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
	fmt.Printf("LeftRotate %d\n", n.data)
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
	fmt.Printf("RightRotate %d\n", n.data)
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
	fmt.Printf("Fixing insertion of %d\n", n.data)
	for n != rbt.root && n.parent.color == Red {
		parent := n.parent
		grandparent := parent.parent
		fmt.Printf("parent %v, grandparent %v\n", parent, grandparent)

		if parent == grandparent.left {
			fmt.Printf("uncle (grandparent.right) %v\n", grandparent.right)
			uncle := grandparent.right
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle != nil && uncle.color == Red {
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
			fmt.Printf("uncle (grandparent.left) %v\n", grandparent.left)
			uncle := grandparent.left
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle != nil && uncle.color == Red {
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
	fmt.Printf("Inserting %d\n", data)
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
	if parent == nil {
		rbt.root = new_node
		new_node.color = Black
		return true
	} else if new_node.data < parent.data {
		parent.left = new_node
	} else {
		parent.right = new_node
	}

	if new_node.parent == nil {
		fmt.Printf("new_node.parent == nil, root data %d\n", new_node.data)
		new_node.color = Black
		return true
	}

	if new_node.parent.parent == nil {
		fmt.Printf("new_node.parent.parent == nil, data: %d\n", new_node.data)
		return true
	}

	rbt.FixInsert(new_node)
	return true
}

// TODO: func (rbt *RBTree) Search(data int) bool {}
func Preorder(n *Node) {
	if n != nil {
		if n.color {
			fmt.Printf("%s%d%s ", colorRed, n.data, colorNone)
		} else {
			fmt.Printf("%d ", n.data)
		}
		Preorder(n.left)
		Preorder(n.right)
	}
}
func Postorder(n *Node) {
	if n != nil {
		Preorder(n.left)
		if n.color {
			fmt.Printf("%s%d%s ", colorRed, n.data, colorNone)
		} else {
			fmt.Printf("%d ", n.data)
		}
		Preorder(n.right)
	}
}

func inorderPrint(n *Node) {
	if n != nil {
		Preorder(n.left)
		Preorder(n.right)
		if n.color {
			fmt.Printf("%s%d%s ", colorRed, n.data, colorNone)
		} else {
			fmt.Printf("%d ", n.data)
		}
	}
}

func Inorder(n *Node) {
	fmt.Print("Inorder: ")
	inorderPrint(n)
	fmt.Println()
}

func columnCount(h int) int {
	if h == 1 {
		return 1
	}
	return columnCount(h-1) + columnCount(h-1) + 1
}

func GetTreeHeight(root *Node) int {
	if root == nil {
		return 0
	}
	return max(GetTreeHeight(root.left), GetTreeHeight(root.right)) + 1
}

func populateTree(M [][]*Node, root *Node, col, row, h int) {
	if root == nil {
		return
	}
	M[row][col] = root
	populateTree(M, root.left, col-int(math.Pow(2, float64(h-2))), row+1, h-1)
	populateTree(M, root.right, col+int(math.Pow(2, float64(h-2))), row+1, h-1)
}

func (tree *RBTree) TreePrinter() {
	height := GetTreeHeight(tree.root)
	columns := columnCount(height)

	matrix := make([][]*Node, height)
	for i := range matrix {
		matrix[i] = make([]*Node, columns)
	}

	populateTree(matrix, tree.root, columns/2, 0, height)

	for i := 0; i < height; i++ {
		for j := 0; j < columns; j++ {
			if matrix[i][j] == nil {
				fmt.Print("  ")
			} else {
				node := matrix[i][j]
				if node.color == Red {
					fmt.Printf("%s%2d%s ", colorRed, node.data, colorNone)
				} else {
					fmt.Printf("%2d ", node.data)
				}
			}
		}
		fmt.Println()
	}
}

func main() {
	rbt := NewRedBlackTree()
	vals := []int{10, 11, 9, 9, 15, 35, 25, 20, 1, 12,
		17, 13, 15, 114, 23, 75, 64, 32, 74, 99, 43}
	for _, val := range vals {
		rbt.Insert(val)
	}
	fmt.Println(GetTreeHeight(rbt.root))
	rbt.TreePrinter()
}
