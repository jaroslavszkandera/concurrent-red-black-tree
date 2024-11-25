package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkRedBlackTree_Insert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		elems := generateRandomNumbers(n)
		rbt := NewRedBlackTree()
		for _, elem := range elems {
			rbt.Insert(elem)
		}
	}
}

func generateRandomNumbers(size int) []int {
	numbers := make([]int, size)
	for i := range numbers {
		numbers[i] = rand.Intn(size * 10)
	}
	return numbers
}

func TestRedBlackTree_Insert(t *testing.T) {
	rbt := NewRedBlackTree()
	rbt.Insert(10)
	if rbt.root.color != Black || rbt.root.data != 10 {
		t.Errorf("Expected root to be black (10), got %v", rbt.root)
	}

	elements := []int{5, 15, 3, 7, 12, 17, 1, 4, 8}
	for _, element := range elements {
		rbt.Insert(element)
	}

	expected := "1 3 4 5 7 8 10 12 15 17 "
	var traversal string
	inorder(rbt.root, func(n *Node) {
		traversal += fmt.Sprintf("%d ", n.data)
	})
	if traversal != expected {
		t.Errorf("Expected inorder traversal: %s, got: %s", expected, traversal)
	}

	checkNoRedChildren(t, rbt.root)
	checkRootBlack(t, rbt.root)
}

func checkNoRedChildren(t *testing.T, node *Node) {
	if node == nil {
		return
	}
	if node.color == Red {
		if node.left != nil && node.left.color == Red {
			t.Errorf("Red node %d has a red left child", node.data)
		}
		if node.right != nil && node.right.color == Red {
			t.Errorf("Red node %d has a red right child", node.data)
		}
	}
	checkNoRedChildren(t, node.left)
	checkNoRedChildren(t, node.right)
}

func checkRootBlack(t *testing.T, node *Node) {
	if node != nil && node.color != Black {
		t.Errorf("Root node is not black")
	}
}

func inorder(n *Node, f func(*Node)) {
	if n != nil {
		inorder(n.left, f)
		f(n)
		inorder(n.right, f)
	}
}
