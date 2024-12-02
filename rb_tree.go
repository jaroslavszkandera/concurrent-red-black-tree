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

// TODO(JS): sentinel Node
type RBTree struct {
	root *Node
}

func NewNode(data int) *Node {
	return &Node{data: data, color: Red, left: nil, right: nil, parent: nil}
}

func NewRedBlackTree() *RBTree {
	return &RBTree{root: nil}
}

func (rbt *RBTree) minimum(n *Node) *Node {
	curr := n
	for curr.left != nil {
		curr = curr.left
	}
	return curr
}

func getSibling(n *Node) *Node {
	if n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func (rbt *RBTree) successor(n *Node) *Node {
	if n == nil {
		return nil
	}
	if n.right != nil {
		return rbt.minimum(n.right)
	}
	m := n.parent
	for m != nil && n == m.right {
		m = n
		m = m.parent
	}
	return m
}

func (rbt *RBTree) Delete(data int) bool {
	fmt.Printf("Deleting %d\n", data)
	if rbt.root == nil {
		return false
	}
	nodeToDel := rbt.Search(data)
	if nodeToDel == nil {
		return false // not found
	}
	fmt.Println("Found nodeToDel", nodeToDel)

	var child *Node
	var replacement *Node
	if nodeToDel.left == nil || nodeToDel.right == nil {
		// has at most one child
		replacement = nodeToDel
	} else {
		// successor
		replacement = rbt.successor(nodeToDel)
	}

	if replacement.left != nil {
		child = replacement.left
	} else {
		child = replacement.right
	}

	fmt.Printf("child %v replacement %v\n", child, replacement)
	if child != nil {
		child.parent = replacement.parent
	}

	// Update parent to point to the child
	if replacement.parent == nil {
		rbt.root = child
	} else if replacement == replacement.parent.left {
		replacement.parent.left = child
	} else {
		replacement.parent.right = child
	}

	if replacement != nodeToDel {
		nodeToDel.data = replacement.data
	}
	// this is not alright, we need the sentinel
	if replacement.color == Black && child != nil {
		rbt.fixDelete(child)
	}

	return true
}

func (rbt *RBTree) fixDelete(n *Node) {
	fmt.Println("fixDelete", n)
	for n != rbt.root && (n == nil || n.color == Black) {
		if n == n.parent.left {
			sibling := n.parent.right
			if sibling.color == Red {
				sibling.color = Black
				n.parent.color = Red
				rbt.leftRotate(n.parent)
				sibling = n.parent.right
			}
			if (sibling.left == nil || sibling.left.color == Black) &&
				(sibling.right == nil || sibling.right.color == Black) {
				sibling.color = Red
				n = n.parent
			} else {
				if sibling.right == nil || sibling.right.color == Black {
					if sibling.left != nil {
						sibling.left.color = Black
					}
					sibling.color = Red
					rbt.rightRotate(sibling)
					sibling = n.parent.right
				}
				sibling.color = n.parent.color
				n.parent.color = Black
				if sibling.right != nil {
					sibling.right.color = Black
				}
				rbt.leftRotate(n.parent)
				n = rbt.root
			}
		} else { // Node is a right child
			sibling := n.parent.left
			if sibling.color == Red {
				sibling.color = Black
				n.parent.color = Red
				rbt.rightRotate(n.parent)
				sibling = n.parent.left
			}
			if (sibling.left == nil || sibling.left.color == Black) &&
				(sibling.right == nil || sibling.right.color == Black) {
				sibling.color = Red
				n = n.parent
			} else {
				if sibling.left == nil || sibling.left.color == Black {
					if sibling.right != nil {
						sibling.right.color = Black
					}
					sibling.color = Red
					rbt.leftRotate(sibling)
					sibling = n.parent.left
				}
				sibling.color = n.parent.color
				n.parent.color = Black
				if sibling.left != nil {
					sibling.left.color = Black
				}
				rbt.rightRotate(n.parent)
				n = rbt.root
			}
		}
	}
	if n != nil {
		n.color = Black
	}
}

func (rbt *RBTree) leftRotate(n *Node) {
	fmt.Printf("leftRotate %d\n", n.data)
	newParent := n.right
	n.right = newParent.left

	if newParent.left != nil {
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

	if newParent.right != nil {
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
					rbt.leftRotate(n)
				}
				// Case 2.2 (Uncle is black and n is the left child):
				//   Recolor parent and grandparent, right rotate on grandparent
				parent.color = Black
				grandparent.color = Red
				rbt.rightRotate(grandparent)
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
					rbt.rightRotate(n)
				}
				parent.color = Black
				grandparent.color = Red
				rbt.leftRotate(grandparent)
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

func (rbt *RBTree) Search(data int) *Node {
	curr := rbt.root
	for curr != nil {
		if data < curr.data {
			curr = curr.left
		} else if data > curr.data {
			curr = curr.right
		} else {
			return curr
		}
	}
	return nil
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
	vals1 := []int{10, 11, 9, 9, 15, 35, 25, 20, 1, 12,
		17, 13, 15, 114, 23, 75, 64, 32, 74, 99, 43}
	vals2 := []int{42, 65, 74, 90, 64, 85, 92, 48}
	for _, val := range vals1 {
		rbt.Insert(val)
	}
	fmt.Println(GetTreeHeight(rbt.root))
	rbt.TreePrinter()
	fmt.Println(rbt.Search(64))
	rbt.Delete(15)
	rbt.TreePrinter()
	rbt.Delete(99)
	rbt.TreePrinter()
	rbt.Delete(35)
	rbt.TreePrinter()
	fmt.Printf("root %v root.left %v root.right %v", rbt.root.data, rbt.root.left.data, rbt.root.right.data)
	for i, val1 := range vals1 {
		rbt.TreePrinter()
		rbt.Delete(val1)
		if i < len(vals2) {
			rbt.Insert(vals2[i])
		}
	}
	rbt.TreePrinter()
}
