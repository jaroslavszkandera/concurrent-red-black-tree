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
	RED   Color = true
	BLACK       = false
)

var NIL *Node = &Node{color: BLACK}

type Node struct {
	color               Color
	data                int
	left, right, parent *Node
}

// TODO(JS): sentinel Node
type RBTree struct {
	root *Node
}

func NewNode(data int) *Node {
	return &Node{data: data, color: RED}
}

func NewRedBlackTree() *RBTree {
	return &RBTree{root: NIL}
}

func (rbt *RBTree) minimum(n *Node) *Node {
	curr := n
	for curr.left != NIL {
		curr = curr.left
	}
	return curr
}

// pointers must be fixed by caller
func (rbt *RBTree) transplant(n, m *Node) {
	if n.parent == nil {
		rbt.root = m
	} else if n == n.parent.left {
		n.parent.left = m
	} else {
		n.parent.right = m
	}
	// if m != NIL {
	m.parent = n.parent
	// }
}

func (rbt *RBTree) Delete(data int) bool {
	fmt.Printf("Deleting %d\n", data)
	nodeToDel := rbt.Search(data)
	if nodeToDel == NIL {
		return false // not found
	}
	fmt.Println("Found nodeToDel", nodeToDel, "root", rbt.root)

	y := nodeToDel
	origColor := y.color
	var x *Node

	if nodeToDel.left == NIL { // Case 1
		x = nodeToDel.right
		rbt.transplant(nodeToDel, nodeToDel.right)
	} else if nodeToDel.right == NIL { // Case 2
		x = nodeToDel.left
		rbt.transplant(nodeToDel, nodeToDel.left)
	} else { // Case 3
		y = rbt.minimum(nodeToDel.right)
		origColor = y.color
		x = y.right

		if y.parent == nodeToDel {
			x.parent = y
		} else {
			rbt.transplant(y, y.right)
			y.right = nodeToDel.right
			y.right.parent = y
		}
		rbt.transplant(nodeToDel, y)
		y.left = nodeToDel.left
		y.left.parent = y
		y.color = nodeToDel.color
	}
	if origColor == BLACK {
		fmt.Println("about to fixDelete", x, y, origColor, nodeToDel)
		rbt.fixDelete(x)
	}
	return true
}

func (rbt *RBTree) fixDelete(n *Node) {
	if n == nil {
		panic("nil fixDelete")
	}
	for n != rbt.root && n.color == BLACK {
		if n == n.parent.left {
			sibling := n.parent.right
			// Type 1
			if sibling.color == RED {
				sibling.color = BLACK
				n.parent.color = RED
				rbt.leftRotate(n.parent)
				sibling = n.parent.right
			}
			// Type 2
			if sibling == NIL || sibling.right == NIL || sibling.left == NIL {
				fmt.Println(n, n.parent, n.parent.right)
				panic("invalid sibling")
			}
			if sibling.left.color == BLACK && sibling.right.color == BLACK {
				sibling.color = RED
				n = n.parent
			} else {
				// Type 3
				if sibling.right.color == BLACK {
					sibling.left.color = BLACK
					sibling.color = RED
					rbt.rightRotate(sibling)
					sibling = n.parent.right
				}
				// Type 4
				sibling.color = n.parent.color
				n.parent.color = BLACK
				sibling.right.color = BLACK
				rbt.leftRotate(n.parent)
				n = rbt.root
			}
		} else {
			sibling := n.parent.left
			// Type 1
			if sibling.color == RED {
				sibling.color = BLACK
				n.parent.color = RED
				rbt.rightRotate(n.parent)
				sibling = n.parent.left
			}
			// Type 2
			if sibling.right.color == BLACK && sibling.left.color == BLACK {
				sibling.color = RED
				n = n.parent
			} else {
				// Type 3
				if sibling.left.color == BLACK {
					sibling.right.color = BLACK
					sibling.color = RED
					rbt.leftRotate(sibling)
					sibling = n.parent.left
				}
				// Type 4
				sibling.color = n.parent.color
				n.parent.color = BLACK
				sibling.left.color = BLACK
				rbt.rightRotate(n.parent)
				n = rbt.root
			}
		}
	}
	n.color = BLACK
}

func (rbt *RBTree) leftRotate(n *Node) {
	fmt.Printf("leftRotate %d\n", n.data)
	newParent := n.right
	n.right = newParent.left

	if newParent.left != NIL {
		newParent.left.parent = n
	}
	newParent.parent = n.parent
	if n.parent == nil {
		rbt.root = newParent
	} else if n == n.parent.left {
		n.parent.left = newParent
	} else {
		n.parent.right = newParent
	}
	newParent.left = n
	n.parent = newParent
}

func (rbt *RBTree) rightRotate(n *Node) {
	fmt.Printf("rightRotate %d\n", n.data)
	newParent := n.left
	n.left = newParent.right

	if newParent.right != NIL {
		newParent.right.parent = n
	}
	newParent.parent = n.parent
	if n.parent == nil {
		rbt.root = newParent
	} else if n == n.parent.right {
		n.parent.right = newParent
	} else {
		n.parent.left = newParent
	}
	newParent.right = n
	n.parent = newParent
}

// Fix red-red situation
func (rbt *RBTree) FixInsert(n *Node) {
	fmt.Printf("Fixing insertion of %d\n", n.data)
	for n.parent != nil && n.parent.color == RED {
		parent := n.parent
		grandparent := parent.parent
		fmt.Printf("parent %v, grandparent %v\n", parent, grandparent)

		if parent == grandparent.left {
			fmt.Printf("uncle (grandparent.right) %v\n", grandparent.right)
			uncle := grandparent.right
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle.color == RED {
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				n = grandparent
			} else {
				// Case 2.1 (Uncle is black and n is the right child):
				//   Perform a left rotation on parent
				if n == parent.right {
					n = parent
					rbt.leftRotate(n)
				}
				// Case 2.2 (Uncle is black and n is the left child):
				//   Recolor parent and grandparent, right rotate on grandparent
				parent.color = BLACK
				grandparent.color = RED
				rbt.rightRotate(grandparent)
			}
		} else {
			fmt.Printf("uncle (grandparent.left) %v\n", grandparent.left)
			uncle := grandparent.left
			// Case 1 (Uncle is red): Recolor parent and uncle to black,
			//   grandparent to red
			if uncle.color == RED {
				parent.color = BLACK
				uncle.color = BLACK
				grandparent.color = RED
				n = grandparent
			} else {
				// Case 2 (Uncle is black):
				if n == parent.left {
					n = parent
					rbt.rightRotate(n)
				}
				parent.color = BLACK
				grandparent.color = RED
				rbt.leftRotate(grandparent)
			}
		}
	}
	rbt.root.color = BLACK
}

func (rbt *RBTree) Insert(data int) bool {
	fmt.Printf("Inserting %d\n", data)
	newNode := &Node{color: RED, data: data, left: NIL, right: NIL}
	var parent *Node
	curr := rbt.root

	// BTS insert
	for curr != NIL {
		parent = curr
		fmt.Println("print", newNode, curr, rbt.root)
		if newNode.data < curr.data {
			curr = curr.left
		} else if newNode.data > curr.data {
			curr = curr.right
		} else {
			return false // Duplicate
		}
	}

	newNode.parent = parent
	if parent == nil {
		rbt.root = newNode
		newNode.color = BLACK
		return true // skip fixInsert
	} else if newNode.data < parent.data {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	rbt.FixInsert(newNode)
	return true
}

func (rbt *RBTree) Search(data int) *Node {
	curr := rbt.root
	for curr != NIL {
		if data < curr.data {
			curr = curr.left
		} else if data > curr.data {
			curr = curr.right
		} else {
			return curr
		}
	}
	return NIL
}

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

func (rbt *RBTree) TreePrinter() {
	if rbt.root == nil {
		fmt.Println("<nil>")
		return
	}
	height := GetTreeHeight(rbt.root)
	columns := columnCount(height)

	matrix := make([][]*Node, height)
	for i := range matrix {
		matrix[i] = make([]*Node, columns)
	}

	populateTree(matrix, rbt.root, columns/2, 0, height)

	for i := 0; i < height; i++ {
		for j := 0; j < columns; j++ {
			if matrix[i][j] == nil {
				fmt.Print("  ")
			} else {
				node := matrix[i][j]
				if node.color == RED {
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
	vals1 := []int{10, 11, 9, 9, 15, 35, 25, 20, 1, 12,
		17, 13, 15, 114, 23, 75, 64, 32, 74, 99, 43}
	// vals2 := []int{42, 65, 74, 90, 64, 85, 92, 48}
	for _, val := range vals1 {
		rbt.Insert(val)
		rbt.TreePrinter()
	}
	fmt.Println(GetTreeHeight(rbt.root))
	rbt.TreePrinter()
	// fmt.Println(rbt.Search(64))
	// rbt.Delete(15)
	// rbt.TreePrinter()
	// rbt.Delete(99)
	// rbt.TreePrinter()
	// rbt.Delete(35)
	// rbt.TreePrinter()
	// fmt.Printf("root %v root.left %v root.right %v", rbt.root.data, rbt.root.left.data, rbt.root.right.data)
	// for i, _ := range vals1 {
	// 	rbt.TreePrinter()
	// 	// rbt.Delete(val1)
	// 	if i < len(vals2) {
	// 		rbt.Insert(vals2[i])
	// 	}
	// }
	// rbt.TreePrinter()
}
