// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper

// BstNode represents the binary search tree node.
type BstNode[T Number] struct {
	value T
	left  *BstNode[T]
	right *BstNode[T]
}

// Bst represents the binary search tree.
type Bst[T Number] struct {
	root *BstNode[T]
}

// NewBst creates a new binary search tree.
func NewBst[T Number]() *Bst[T] {
	return &Bst[T]{}
}

// Insert adds a new value to the binary search tree.
func (b *Bst[T]) Insert(value T) {
	node := &BstNode[T]{
		value: value,
	}

	if b.root == nil {
		b.root = node
		return
	}

	cur := b.root

	for {
		if node.value <= cur.value {
			if cur.left == nil {
				cur.left = node
				return
			}

			cur = cur.left
		} else {
			if cur.right == nil {
				cur.right = node
				return
			}

			cur = cur.right
		}
	}
}

// Contains checks whether the given value exists in the binary search tree.
func (b *Bst[T]) Contains(value T) bool {
	node, _ := b.searchNode(value)
	return node != nil
}

// Remove removes the specified value from the binary search tree
// and rebalances the tree.
func (b *Bst[T]) Remove(value T) bool {
	node, parent := b.searchNode(value)
	if node == nil {
		return false
	}

	b.removeNode(node, parent)
	return true
}

// Min function returns the minimum value in the binary search tree.
func (b *Bst[T]) Min() T {
	if b.root == nil {
		return T(0)
	}

	node, _ := minNode(b.root)
	return node.value
}

// Max function returns the maximum value in the binary search tree.
func (b *Bst[T]) Max() T {
	if b.root == nil {
		return T(0)
	}

	node, _ := maxNode(b.root)
	return node.value
}

// searchNode searches for the given value in the binary search tree and returns
// the first matching node and its parent.
func (b *Bst[T]) searchNode(value T) (*BstNode[T], *BstNode[T]) {
	var parent *BstNode[T]
	node := b.root

	for node != nil {
		diff := value - node.value
		if diff == 0 {
			break
		}

		parent = node
		if diff < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node, parent
}

// removeNode removes the specified node from the binary search tree
// and rebalances the tree.
func (b *Bst[T]) removeNode(node, parent *BstNode[T]) {
	if node.left != nil && node.right != nil {
		min, minParent := minNode(node.right)
		if minParent == nil {
			minParent = node
		}

		b.removeNode(min, minParent)
		node.value = min.value
	} else {
		var child *BstNode[T]

		if node.left != nil {
			child = node.left
		} else {
			child = node.right
		}

		if node == b.root {
			b.root = child
		} else if parent.left == node {
			parent.left = child
		} else {
			parent.right = child
		}
	}
}

// minNode functions returns the node with the minimum value and its parent node.
func minNode[T Number](root *BstNode[T]) (*BstNode[T], *BstNode[T]) {
	var parent *BstNode[T]
	node := root

	for node.left != nil {
		parent = node
		node = node.left
	}

	return node, parent
}

// maxNode functions returns the node with the maximum value and its parent node.
func maxNode[T Number](root *BstNode[T]) (*BstNode[T], *BstNode[T]) {
	var parent *BstNode[T]
	node := root

	for node.right != nil {
		parent = node
		node = node.right
	}

	return node, parent
}
