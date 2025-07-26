package rbtree

import (
	"bytes"
	"sync"
)

type TreeOrientation int

const (
	MinFirst TreeOrientation = iota
	MaxFirst
)

// A RedBlack tree must follow the following 4 rules:
// 0 - each node must be either red or black
// 1 - root must always be black
// 2 - a red node must have a black parent and black children
// 3 - every path from root to nils passes << exactly >> the same number of black nodes
type Tree struct {
	// The root of the tree
	Root *Node

	// A pool to optimize memory allocations for nodes
	pool *sync.Pool

	// Whether to point the head of the tree to its minimum or maximum node.
	orientation TreeOrientation

	// A pointer to the node with the minimum or maximum price depending on the
	// tree's orientation.
	// Reads in O(1). Refreshes in O(log n)
	head *Node
}

// PrintTree prints the tree in a string that can be easily visualized in the console
func (t *Tree) PrintTree() string {

	if t.Root != nil {
		var buffer bytes.Buffer
		t.Root.PrintTree(&buffer, []byte{}, []byte{})
		return buffer.String()
	}

	return "nil"
}

func NewTree(orientation TreeOrientation, pool *sync.Pool) *Tree {
	t := &Tree{
		orientation: orientation,
		pool:        pool,
	}
	return t
}

func (t *Tree) getNode() *Node {
	n := t.pool.Get().(*Node)
	// Reset the node in case it was not clean
	n.Reset()
	return n
}

func (t *Tree) putNode(n *Node) {
	// Reset the node to push a clean object
	n.Reset()
	t.pool.Put(n)
}
