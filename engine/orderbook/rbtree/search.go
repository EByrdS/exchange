package rbtree

import (
	"fmt"
)

// Min returns the tree's minimum node
//
// O(log n)
func (t *Tree) Min() *Node {
	return t.min(t.Root)
}

// Max returns the tree's maximum node
//
// O(log n)
func (t *Tree) Max() *Node {
	return t.max(t.Root)
}

// Head returns the minimum or maximum node depending on the tree's orientation.
//
// O(1)
func (t *Tree) Head() *Node {
	return t.head
}

// Orientation returns whether this tree is keeping track of the minimum or
// the maximum node.
func (t *Tree) Orientation() TreeOrientation {
	return t.orientation
}

// Get returns the node containing the requested price if it exists
//
// O(log n)
func (t *Tree) Get(price uint64) (*Node, error) {
	node, _, _ := t.find(price)
	if node == nil {
		return nil, fmt.Errorf("price %v does not exist in tree", price)
	}

	return node, nil
}

func (t *Tree) min(min *Node) *Node {
	if min == nil {
		return nil
	}

	for {
		if min.Left == nil {
			return min
		}
		min = min.Left
	}
}

func (t *Tree) max(max *Node) *Node {
	if max == nil {
		return nil
	}

	for {
		if max.Right == nil {
			return max
		}
		max = max.Right
	}
}

// Find receives the target price and returns a reference to the node that
// holds said price if it exists, whether said node is the left or right child,
// and a reference to its parent. This information can be used to insert a new
// node immediately if needed.
//
// If the holder is nil, then the searched element is not in the tree.
// If the parent is nil, then the returned element is the root.
func (t *Tree) find(price uint64) (holder *Node, side ChildSide, parent *Node) {
	if t.Root == nil {
		return nil, LeftChild, nil
	}

	if t.Root.Price == price {
		return t.Root, LeftChild, nil
	}

	parent = t.Root
	for {
		if price > parent.Price {
			if parent.Right == nil {
				return nil, RightChild, parent
			}

			if parent.Right.Price == price {
				return parent.Right, RightChild, parent
			}

			parent = parent.Right
		} else {
			if parent.Left == nil {
				return nil, LeftChild, parent
			}

			if parent.Left.Price == price {
				return parent.Left, LeftChild, parent
			}

			parent = parent.Left
		}
	}
}
