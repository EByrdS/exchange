package rbtree

import (
	"bytes"
	"exchange/engine/orderbook/pricelevel"
	"strconv"
)

// ChildSide is a type to identify the direction of a child
type ChildSide int

const (
	LeftChild ChildSide = iota
	RightChild
)

// Node is the building block of a Red-Black tree, and contains a pointer to a piece of
// data of generic type.
//
// A nil node is evaluated as being black.
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Red    bool

	Price  uint64
	Orders *pricelevel.PriceLevel
}

func NewNode() *Node {
	return &Node{
		Orders: pricelevel.New(),
	}
}

// PrintTree uses recursion to create a text representation of the Node's subtree.
// To get the string of an entire tree, call PrintTree on its node.
func (n *Node) PrintTree(buffer *bytes.Buffer, prefix []byte, childrenPrefix []byte) {
	color := []byte("BLK")
	if n.Red {
		color = []byte("RED")
	}

	buffer.Write(prefix)
	buffer.Write(color)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatUint(n.Price, 10))
	buffer.WriteByte('\n')

	if n.Right != nil {

		rightLink := "->"
		if n.Right.Parent == nil {
			rightLink = "o>"
		} else if n.Right.Parent != n {
			rightLink = "?>"
		}

		char := "├"
		if n.Left == nil {
			char = "└"
		}

		leftChildPrefix := "|"
		if n.Left == nil {
			leftChildPrefix = " "
		}

		var pr []byte
		pr = make([]byte, 0, 20)
		pr = append(pr, childrenPrefix...)
		pr = append(pr, char...)
		pr = append(pr, "R-"...)
		pr = append(pr, rightLink...)

		// var
		var cpr []byte
		cpr = make([]byte, 0, 20)
		cpr = append(cpr, childrenPrefix...)
		cpr = append(cpr, leftChildPrefix...)
		cpr = append(cpr, "   "...)
		n.Right.PrintTree(buffer, pr, cpr)
	}

	if n.Left != nil {

		leftLink := "->"
		if n.Left.Parent == nil {
			leftLink = "o>"
		} else if n.Left.Parent != n {
			leftLink = "?>"
		}

		var pl []byte
		pl = make([]byte, 0, 20)
		pl = append(pl, childrenPrefix...)
		pl = append(pl, "└L-"...)
		pl = append(pl, leftLink...)

		var cpl []byte
		cpl = make([]byte, 0, 20)
		cpl = append(cpl, childrenPrefix...)
		cpl = append(cpl, "    "...)
		n.Left.PrintTree(buffer, pl, cpl)
	}

}

func (n *Node) IsRoot() bool {
	return n.Parent == nil
}

func (n *Node) IsRed() bool {
	return n != nil && n.Red
}

func (n *Node) IsBlack() bool {
	return n == nil || !n.Red
}

func (n *Node) Reset() {
	if n == nil {
		return
	}

	n.Parent = nil
	n.Left = nil
	n.Right = nil
	n.Red = false

	n.Price = 0
	n.Orders.Reset()

	// n.orders.Init()
}
