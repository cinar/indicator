package bst

import (
	"github.com/cinar/indicator/container"
)

// BST node.
type Node struct {
	value interface{}
	left  *Node
	right *Node
}

// BST type.
type Tree struct {
	root *Node
}

// New binary search tree.
func New() *Tree {
	return &Tree{}
}

// Inserts the given value.
func (t *Tree) Insert(value interface{}) {
	newNode := &Node{
		value: value,
	}

	if t.root == nil {
		t.root = newNode
		return
	}

	curNode := t.root

	for {
		if container.Compare(newNode.value, curNode.value) <= 0 {
			if curNode.left == nil {
				curNode.left = newNode
				return
			} else {
				curNode = curNode.left
			}
		} else {
			if curNode.right == nil {
				curNode.right = newNode
				return
			} else {
				curNode = curNode.right
			}
		}
	}
}

// Removes the given value.
func (t *Tree) Remove(value interface{}) bool {
	var parent *Node
	node := t.root

	for node != nil {
		switch container.Compare(value, node.value) {
		case 0:
			t.removeNode(parent, node)
			return true

		case -1:
			parent = node
			node = node.left

		case 1:
			parent = node
			node = node.right
		}
	}

	return false
}

// Min value.
func (t *Tree) Min() interface{} {
	node, _ := minNode(t.root)
	if node == nil {
		return nil
	}

	return node.value
}

// Max value.
func (t *Tree) Max() interface{} {
	node, _ := maxNode(t.root)
	if node == nil {
		return nil
	}

	return node.value
}

// Remove node.
func (t *Tree) removeNode(parent, node *Node) {
	if node.left != nil && node.right != nil {
		min, minParent := minNode(node.right)
		if minParent == nil {
			minParent = node
		}

		t.removeNode(minParent, min)
		node.value = min.value
	} else {
		var child *Node
		if node.left != nil {
			child = node.left
		} else {
			child = node.right
		}

		if node == t.root {
			t.root = child
		} else if parent.left == node {
			parent.left = child
		} else {
			parent.right = child
		}
	}
}

// Min node. Returns min node and its parent.
func minNode(root *Node) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}

	var parent *Node
	node := root

	for node.left != nil {
		parent = node
		node = node.left
	}

	return node, parent
}

// Max node. Returns max node and its parent.
func maxNode(root *Node) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}

	var parent *Node
	node := root

	for node.right != nil {
		parent = node
		node = node.right
	}

	return node, parent
}
